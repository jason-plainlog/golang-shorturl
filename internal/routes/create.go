package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

// Create shorturl (id, url, expireAt) record from request
func Create(c echo.Context) error {
	req := new(request)

	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	res := response{}

	return c.JSON(http.StatusOK, res)
}
