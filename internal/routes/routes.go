package routes

import (
	"github.com/gin-contrib/cors"
	"os"
	"path/filepath"
	"prin/internal/app"
	"prin/internal/routes/qrcode"
)

func RegisterStatic() {
	r := app.App.Gin

	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := filepath.Dir(pwd)

	r.StaticFile("/", path+"/public/index.html")

	for _, v := range []string{"favicon.ico"} {
		r.StaticFile(v, path+"/public/"+v)
	}

	for _, v := range []string{"static"} {
		r.Static("/"+v, path+"/public/"+v)
	}
}

func RegisterRoutes() {
	r := app.App.Gin

	c := cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	r.Use(c)

	// 前端路由
	api := r.Group("/api")
	{
		api.POST("/qrcode", qrcode.Generate)
	}

	RegisterStatic()
}
