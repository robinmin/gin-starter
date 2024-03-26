package middleware_test

// func setup() *middleware.JWTTokenPair {
// 	redisStore := persistence.NewRedisCache("localhost:6379", "", 10*time.Minute)
// 	return &middleware.JWTTokenPair{
// 		AccessSecret:    "test_access_secret",
// 		RefreshSecret:   "test_refresh_secret",
// 		RefreshDelay:    120,
// 		ApplicationName: "TestApp",
// 		RDB:             redisStore,
// 	}
// }

// func TestGenerateTokenPair(t *testing.T) {
// 	testCases := []struct {
// 		username        string
// 		defaultDuration int
// 		expectError     bool
// 		errorInfo       string
// 	}{
// 		{"testuser1", 10, false, "Failed for testuser1 with duration 10"},
// 		{"testuser2", 5, false, "Failed for testuser2 with duration 5"},
// 		// 添加更多测试案例
// 	}

// 	for _, tc := range testCases {
// 		jtp := setup()
// 		tokenPair, err := jtp.GenerateTokenPair(tc.username, tc.defaultDuration)
// 		if tc.expectError {
// 			assert.Error(t, err, tc.errorInfo)
// 		} else {
// 			assert.NoError(t, err, tc.errorInfo)
// 			assert.NotNil(t, tokenPair, tc.errorInfo)
// 			assert.NotEmpty(t, tokenPair.AccessToken, tc.errorInfo)
// 			assert.NotEmpty(t, tokenPair.RefreshToken, tc.errorInfo)
// 		}
// 	}
// }

// func TestGenerateToken(t *testing.T) {
// 	testCases := []struct {
// 		username       string
// 		secret         string
// 		expiryDuration int
// 		expectError    bool
// 		errorInfo      string
// 	}{
// 		{"user1", "secret1", 10, false, "GenerateToken should succeed for user1"},
// 		// 添加更多测试案例
// 	}

// 	for _, tc := range testCases {
// 		jtp := setup()
// 		token, err := jtp.GenerateToken(tc.username, tc.secret, tc.expiryDuration)
// 		if tc.expectError {
// 			assert.Error(t, err, tc.errorInfo)
// 		} else {
// 			assert.NoError(t, err, tc.errorInfo)
// 			assert.NotEmpty(t, token, tc.errorInfo)
// 		}
// 	}
// }

// func TestParseToken(t *testing.T) {
// 	// 注意: 这里可能需要生成一个有效的 token 用于测试
// 	testCases := []struct {
// 		token       string
// 		secret      string
// 		expectError bool
// 		errorInfo   string
// 	}{
// 		{"validToken", "secret", false, "ParseToken should succeed with valid token"},
// 		// 添加更多测试案例
// 	}

// 	for _, tc := range testCases {
// 		jtp := setup()
// 		claims, err := jtp.ParseToken(tc.token, tc.secret)
// 		if tc.expectError {
// 			assert.Error(t, err, tc.errorInfo)
// 		} else {
// 			assert.NoError(t, err, tc.errorInfo)
// 			assert.NotNil(t, claims, tc.errorInfo)
// 		}
// 	}
// }

// func TestIsValidAccessToken(t *testing.T) {
// 	// 这里需要预先在 Redis 中设置一些 token 用于测试
// 	testCases := []struct {
// 		username    string
// 		accessToken string
// 		expectValid bool
// 		errorInfo   string
// 	}{
// 		{"user1", "existingToken", true, "IsValidAccessToken should return true for valid token"},
// 		// 添加更多测试案例
// 	}

// 	for _, tc := range testCases {
// 		jtp := setup()
// 		isValid := jtp.IsValidAccessToken(tc.username, tc.accessToken)
// 		assert.Equal(t, tc.expectValid, isValid, tc.errorInfo)
// 	}
// }

// func TestRefreshTokenPair(t *testing.T) {
// 	// 这里可能需要预先生成并设置一些 token
// 	testCases := []struct {
// 		refreshToken    string
// 		defaultDuration int
// 		expectError     bool
// 		errorInfo       string
// 	}{
// 		{"existingRefreshToken", 10, false, "RefreshTokenPair should succeed with valid refresh token"},
// 		// 添加更多测试案例
// 	}

// 	for _, tc := range testCases {
// 		jtp := setup()
// 		newAccessToken, err := jtp.RefreshTokenPair(tc.refreshToken, tc.defaultDuration)
// 		if tc.expectError {
// 			assert.Error(t, err, tc.errorInfo)
// 		} else {
// 			assert.NoError(t, err, tc.errorInfo)
// 			assert.NotEmpty(t, newAccessToken, tc.errorInfo)
// 		}
// 	}
// }

// func TestReleaseTokenPair(t *testing.T) {
// 	// 这里可能需要预先生成并设置一些 token
// 	testCases := []struct {
// 		accessToken   string
// 		refreshToken  string
// 		expectSuccess bool
// 		errorInfo     string
// 	}{
// 		{"existingAccessToken", "existingRefreshToken", true, "ReleaseTokenPair should succeed with valid tokens"},
// 		// 添加更多测试案例
// 	}

// 	for _, tc := range testCases {
// 		jtp := setup()
// 		success, err := jtp.ReleaseTokenPair(tc.accessToken, tc.refreshToken)
// 		if !tc.expectSuccess {
// 			assert.Error(t, err, tc.errorInfo)
// 		} else {
// 			assert.NoError(t, err, tc.errorInfo)
// 			assert.True(t, success, tc.errorInfo)
// 		}
// 	}
// }
