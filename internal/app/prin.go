package app

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strings"
)

type Instance struct {
	Gin *gin.Engine
}

var (
	// App 程序实例
	App      *Instance
	RootPath string
)

func init() {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "production" || env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化环境配置
	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	RootPath = filepath.Dir(pwd)

	App = &Instance{
		Gin: gin.New(),
	}
}
