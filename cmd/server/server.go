package main

import (
	"context"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/models"
	"url-shortener/internal/routes"
	"url-shortener/internal/token"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	db := mongoClient.Database(cfg.DATABASE)

	// setting up go routine to generate tokens offline
	tokenChan := make(chan string, cfg.MAX_TOKEN)
	go token.GenToken(db.Collection(cfg.RECORD_COLLECTION), tokenChan)

	// setting up memcached cache
	cache := memcache.New(cfg.MEMCACHED_ADDRS...)
	if err := cache.Ping(); err != nil {
		panic(err)
	}

	// Setting up API backend server
	server := echo.New()
	server.Use(middleware.BodyLimit("4K"))

	server.GET("/:id", routes.Get(db, cache))
	server.POST("/api/v1/urls", routes.Create(db, tokenChan))

	server.Logger.Fatal(
		server.Start(cfg.LISTEN_ADDR),
	)
}
