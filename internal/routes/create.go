package routes

import (
	"fmt"
	"net/http"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/models"

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
	if r.URL == "" {
		return fmt.Errorf("url can't be blank")
	}

	if len(r.URL) > cfg.MAX_URL_LEN {
		return fmt.Errorf("the length of url can't exceed %d", cfg.MAX_URL_LEN)
	}

	if time.Now().After(r.ExpireAt) {
		return fmt.Errorf("expireAt should be in the future")
	}

	return nil
}

var cfg = config.GetConfig()

// Create shorturl (id, url, expireAt) record from request
func Create(db *mongo.Database, tokenChan chan string) echo.HandlerFunc {
	records := db.Collection("records")

	return func(c echo.Context) error {
		req := new(request)

		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// check for request validity
		if err := req.CheckValidity(); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		record := models.Record{
			ID:       <-tokenChan,
			URL:      req.URL,
			ExpireAt: req.ExpireAt,
		}

		for err := record.Save(records); err == models.ErrIdInUse; {
			record.ID = <-tokenChan
			err = record.Save(records)
		}

		return c.JSON(http.StatusOK, response{
			ID:       record.ID,
			ShortURL: cfg.BASE_URL + "/" + record.ID,
		})
	}
}
