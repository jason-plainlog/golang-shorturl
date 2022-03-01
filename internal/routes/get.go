package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Resovle shorturl id to redirect url or return 404 error
func Get(c echo.Context) error {
	notFound := true
	expired := false

	if notFound || expired {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "")
}
