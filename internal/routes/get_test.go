package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/models"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo/v4"
)

func TestGetRoute(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	client, _ := models.Connect()

	cache := memcache.New(cfg.MEMCACHED_ADDRS...)
	handler := Get(client.Database("shorturl_test"), cache)

	if err := handler(c); err != nil {
		fmt.Print(err)
		t.Fail()
	}

	if rec.Result().StatusCode != http.StatusNotFound {
		fmt.Print(rec.Result().Status)
		t.Fail()
	}
}
