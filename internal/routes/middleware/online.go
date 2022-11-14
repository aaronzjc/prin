package middleware

import (
	"fmt"
	"net/http"
	"os"
	"prin/internal/util/tool"

	"github.com/gin-gonic/gin"
)

func SetOnline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer ctx.Next()
		clientIp := tool.ClientIp(ctx.Request)
		if clientIp == "" {
			return
		}
		online := os.Getenv("ONLINE_SVC")
		if online == "" {
			return
		}
		url := fmt.Sprintf("%s/online/%s/%s", online, "prin", clientIp)
		req, _ := http.NewRequest("POST", url, nil)
		resp, err := (new(http.Client)).Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}
}
