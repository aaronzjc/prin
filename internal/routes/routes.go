package routes

import (
	"prin/internal/app"
	"prin/internal/routes/handler"
	"prin/internal/routes/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterStatic() {
	r := app.App.Gin
	path := app.RootPath

	r.StaticFile("/", path+"/public/index.html")
	r.StaticFile("/favicon.ico", path+"/public/favicon.ico")

	for k, v := range map[string]string{"static": "/public/static", "assets": "/assets"} {
		r.Static("/"+k, path+v)
	}

	// 自动下载CA证书
	r.GET("/ca.pem", func(c *gin.Context) {
		fileName := "ca.pem"
		target := path + "/assets/ca.pem"
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	})
}

func RegisterRoutes() {
	r := app.App.Gin
	r.Use(gin.Recovery(), gin.Logger())

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
		// 二维码
		qrcode := &handler.Qrcode{}
		api.POST("/qrcode", qrcode.Generate)

		// 编解码
		coder := &handler.Coder{}
		api.POST("/coder", coder.Decode)

		// 证书
		cert := &handler.Cert{}
		api.POST("/cert", cert.Generate)

		// iptables
		iptables := &handler.Iptables{}
		api.POST("/iptables", iptables.Beauty)

		// 统计
		statGroup := api.Group("/stat").Use(middleware.SetOnline())
		{
			stat := &handler.Stat{}
			statGroup.GET("/online", stat.Online)
		}
	}

	RegisterStatic()
}
