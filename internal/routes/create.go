package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	request struct {
		URL      string    `json:"url"`
		ExpireAt time.Time `json:"expireAt"`
	}

	response struct {
		ID       string `json:"id"`
		ShortURL string `json:"shortUrl"`
	}
)

// Check for request validity
func (r *request) CheckValidity() error {
	if time.Now().After(r.ExpireAt) {
		return fmt.Errorf("expireAt should be in the future")
	}

	return nil
}

// Create shorturl (id, url, expireAt) record from request
func Create(db *mongo.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(request)

		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// check for request validity
		if err := req.CheckValidity(); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		res := response{}

		return c.JSON(http.StatusOK, res)
	}
}
