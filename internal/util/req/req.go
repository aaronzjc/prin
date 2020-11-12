package req

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CodeSuccess    = 10000
	CodeError      = 10001
	CodeForbidden  = 10002
	CodeAuthFailed = 10003
)

func JSON(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
