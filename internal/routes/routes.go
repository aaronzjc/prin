package routes

import (
	"github.com/gin-contrib/cors"
	"prin/internal/app"
	"prin/internal/routes/cert"
	"prin/internal/routes/coder"
	"prin/internal/routes/qrcode"
)

func RegisterStatic() {
	r := app.App.Gin
	path := app.RootPath

	r.StaticFile("/", path+"/public/index.html")
	r.StaticFile("/favicon.ico", path+"/public/favicon.ico")

	for k, v := range map[string]string{"static": "/public/static", "assets": "/assets"} {
		r.Static("/"+k, path+v)
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
		api.POST("/coder", coder.Decode)
		api.POST("/cert", cert.Generate)
	}

	RegisterStatic()
}
