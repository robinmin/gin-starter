package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 表示统一响应的JSON格式
type Result struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

func NewResult(code int, message string, data interface{}) Result {
	return Result{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewQuickResult(code int, data interface{}) Result {
	return Result{
		Code:    code,
		Message: GetMessage(code),
		Data:    data,
	}
}

// 接口执行正常 需要返回数据 data
func (result Result) OK(c *gin.Context) {
	c.JSON(http.StatusOK, result)
}

func (result Result) Fail(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, result)
	c.Abort()
}

func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//先执行请求
		c.Next()

		// 发生了错误
		if len(c.Errors) > 0 {
			//获取最后一个error 返回
			err := c.Errors.Last()
			NewResult(http.StatusInternalServerError, err.Error(), nil).Fail(c)
			return
		}
	}
}
