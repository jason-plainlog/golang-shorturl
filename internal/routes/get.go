package routes

import (
	"context"
	"net/http"
	"time"
	"url-shortener/internal/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Resovle shorturl id to redirect url or return 404 error
func Get(db *mongo.Database) echo.HandlerFunc {
	records := db.Collection(cfg.RECORD_COLLECTION)

	return func(c echo.Context) error {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		found := true

		var result models.Record
		if err := records.FindOne(ctx, bson.M{"id": id}).Decode(&result); err != nil {
			if err == mongo.ErrNoDocuments {
				found = false
			} else {
				c.Logger().Fatal(err)
				return c.JSON(http.StatusInternalServerError, nil)
			}
		}

		if !found || result.ExpireAt.Before(time.Now()) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.Redirect(http.StatusTemporaryRedirect, result.URL)
	}
}
