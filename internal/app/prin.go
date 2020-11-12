package app

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

type Instance struct {
	Gin *gin.Engine
}

var (
	// App 程序实例
	App *Instance
)

func init() {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "production" || env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	App = &Instance{
		Gin: gin.New(),
	}
}
