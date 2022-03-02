package routes

import (
	"context"
	"net/http"
	"time"
	"url-shortener/internal/models"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Resovle shorturl id to redirect url or return 404 error
func Get(db *mongo.Database, cache *memcache.Client) echo.HandlerFunc {
	records := db.Collection(cfg.RECORD_COLLECTION)

	return func(c echo.Context) error {
		id := c.Param("id")

		// find record with id in cache first, if cached, redirect immediately
		hit, err := cache.Get(id)
		if err == nil {
			return c.Redirect(http.StatusTemporaryRedirect, string(hit.Value))
		}

		// cache missed, find record with id in mongodb and set cache
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		found := true

		var result models.Record
		if err := records.FindOne(ctx, bson.M{"id": id}).Decode(&result); err != nil {
			if err == mongo.ErrNoDocuments {
				found = false
			} else {
				return err
			}
		}

		if !found || result.ExpireAt.Before(time.Now()) {
			return c.JSON(http.StatusNotFound, nil)
		}

		c.Redirect(http.StatusTemporaryRedirect, result.URL)

		// set cache

		err = cache.Set(&memcache.Item{
			Key:   string(id),
			Value: []byte(result.URL),
		})

		return err
	}
}
