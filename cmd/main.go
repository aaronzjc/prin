package main

import (
	"prin/internal/app"
	"prin/internal/routes"
	"prin/internal/util/logger"
)

func main() {
	routes.RegisterRoutes()

	logger.Fatal(app.App.Gin.Run(":8980"))
}
