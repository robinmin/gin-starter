package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

const (
	AccessTokenPrefix  = "sys_access_"
	RefreshTokenPrefix = "sys_refresh_"
)

type JWTTokenPair struct {
	AccessSecret    string
	RefreshSecret   string
	RefreshDelay    int
	ApplicationName string
	RDB             *persistence.RedisStore
}

type JWTTokenResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token with specific claims.
func (jtp *JWTTokenPair) GenerateToken(username string, secret string, expiryDuration int) (string, error) {
	expireTime := time.Now().Add(time.Duration(expiryDuration) * time.Minute)
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    jtp.ApplicationName,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// GenerateTokenPair generates a pair of access and refresh tokens.
func (jtp *JWTTokenPair) GenerateTokenPair(username string, defaultDuration int) (tkrst *JWTTokenResult, err error) {
	accessToken, err := jtp.GenerateToken(username, jtp.AccessSecret, defaultDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jtp.GenerateToken(username, jtp.RefreshSecret, defaultDuration+jtp.RefreshDelay)
	if err != nil {
		return nil, err
	}

	err = jtp.RDB.Set(AccessTokenPrefix+accessToken, username, time.Duration(defaultDuration)*time.Minute)
	if err != nil {
		return nil, err
	}

	err = jtp.RDB.Set(RefreshTokenPrefix+refreshToken, accessToken, time.Duration(defaultDuration+jtp.RefreshDelay)*time.Minute)
	if err != nil {
		return nil, err
	}

	tkrst = &JWTTokenResult{AccessToken: accessToken, RefreshToken: refreshToken}
	return tkrst, nil
}

// ParseToken parses a JWT token and returns its claims.
func (jtp *JWTTokenPair) ParseToken(token string, secret string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, err
}

// IsValidAccessToken checks if the provided access token is valid and not expired.
func (jtp *JWTTokenPair) IsValidAccessToken(username string, accessToken string) bool {
	var storedUsername string
	err := jtp.RDB.Get(AccessTokenPrefix+accessToken, &storedUsername)
	return err == nil && storedUsername == username
}

// JWTAuthMiddleware creates a middleware for validating access tokens.
func JWTAuthMiddleware(jtp *JWTTokenPair) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未提供访问令牌"})
			ctx.Abort()
			return
		}

		claims, err := jtp.ParseToken(token, jtp.AccessSecret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "无效的令牌"})
			ctx.Abort()
			return
		}

		if !jtp.IsValidAccessToken(claims.Username, token) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "无效或过期的访问令牌"})
			ctx.Abort()
			return
		}

		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}

// RefreshTokenPair uses a refresh token to generate a new access token.
func (jtp *JWTTokenPair) RefreshTokenPair(refreshToken string, defaultDuration int) (newAccessToken string, err error) {
	var oldAccessToken string
	err = jtp.RDB.Get(RefreshTokenPrefix+refreshToken, &oldAccessToken)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	claims, err := jtp.ParseToken(oldAccessToken, jtp.AccessSecret)
	if err != nil {
		return "", errors.New("failed to parse old access token")
	}

	newAccessToken, err = jtp.GenerateToken(claims.Username, jtp.AccessSecret, defaultDuration)
	if err != nil {
		return "", err
	}

	// 设置新的 access token 到 Redis
	err = jtp.RDB.Set(AccessTokenPrefix+newAccessToken, claims.Username, time.Duration(defaultDuration)*time.Minute)
	if err != nil {
		return "", err
	}

	// 可选：更新 refresh token 的过期时间
	err = jtp.RDB.Set(RefreshTokenPrefix+refreshToken, newAccessToken, time.Duration(defaultDuration+jtp.RefreshDelay)*time.Minute)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

// ReleaseTokenPair removes the access and refresh tokens from Redis.
func (jtp *JWTTokenPair) ReleaseTokenPair(accessToken string, refreshToken string) (bool, error) {
	// 删除与 access token 相关联的条目
	err := jtp.RDB.Delete(AccessTokenPrefix + accessToken)
	if err != nil {
		return false, err
	}

	// 删除与 refresh token 相关联的条目
	err = jtp.RDB.Delete(RefreshTokenPrefix + refreshToken)
	if err != nil {
		return false, err
	}

	return true, nil
}

// // generateRandomString generates a random string of specified length.
// func generateRandomString(n int) (string, error) {
// 	bytes := make([]byte, n)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(bytes), nil
// }
