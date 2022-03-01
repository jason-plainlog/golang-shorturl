package models

import (
	"context"
	"time"
	"url-shortener/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cfg = config.GetConfig()

func Connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return mongo.Connect(ctx, options.Client().ApplyURI(cfg.MONGODB_URI))
}
