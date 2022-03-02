package main

import (
	"context"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/models"
	"url-shortener/internal/routes"
	"url-shortener/internal/token"

	"github.com/labstack/echo/v4"
)

var cfg = config.GetConfig()

func main() {
	// Setting up connection to MongoDB
	mongoClient, err := models.Connect()
	if err != nil {
		panic(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	db := mongoClient.Database("shorturl")

	// setting up go routine to generate tokens offline
	tokenChan := make(chan string, cfg.MAX_TOKEN)
	go token.GenToken(db.Collection("records"), tokenChan)

	// Setting up API backend server
	server := echo.New()

	server.GET("/:id", routes.Get(db))
	server.POST("/api/v1/urls", routes.Create(db, tokenChan))

	server.Logger.Fatal(
		server.Start(cfg.LISTEN_ADDR),
	)
}
