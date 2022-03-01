package main

import (
	"url-shortener/internal/config"
	"url-shortener/internal/routes"

	"github.com/labstack/echo/v4"
)

var cfg = config.GetConfig()

func main() {
	server := echo.New()

	server.GET("/:id", routes.Get)
	server.POST("/api/v1/urls", routes.Create)

	server.Logger.Fatal(
		server.Start(cfg.LISTEN_ADDR),
	)
}
