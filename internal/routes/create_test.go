package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/token"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo/v4"
)

func TestCreateRoute(t *testing.T) {
	e := echo.New()
	body, _ := json.Marshal(request{
		URL:      "https://www.dcard.tw/f",
		ExpireAt: time.Now().Add(time.Second * 20),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	client, _ := models.Connect()
	records := client.Database("shorturl_test").Collection(cfg.RECORD_COLLECTION)
	tokenChan := make(chan string, 10)
	go token.GenToken(records, tokenChan)

	handler := Create(client.Database("shorturl_test"), tokenChan)
	if err := handler(c); err != nil {
		t.FailNow()
	}

	if rec.Result().StatusCode != http.StatusOK {
		t.Fail()
	}

	var resp response
	json.Unmarshal(rec.Body.Bytes(), &resp)

	// test get and redirect
	cache := memcache.New(cfg.MEMCACHED_ADDRS...)
	handler = Get(client.Database("shorturl_test"), cache)
	if err := handler(c); err != nil {
		t.FailNow()
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(resp.ID)

	if err := handler(c); err != nil {
		t.FailNow()
	}

	if rec.Result().StatusCode != http.StatusTemporaryRedirect {
		fmt.Print(rec.Result().Status)
		t.Fail()
	}

	url, err := rec.Result().Location()
	if err != nil {
		t.Fail()
	}

	if url.String() != "https://www.dcard.tw/f" {
		t.Fail()
	}
}

func TestRequestValidityCheck(t *testing.T) {
	// current time - 10 second
	// the req should not be valid for expireAt is before current time
	req := request{
		URL:      "http://localhost",
		ExpireAt: time.Now().Add(-time.Second * 10),
	}
	if req.CheckValidity() == nil {
		t.Fail()
	}

	// current time + 10 second, should pass
	req.ExpireAt = time.Now().Add(time.Second * 10)
	if req.CheckValidity() != nil {
		t.Fail()
	}

	// empty url should fail the check too
	req.URL = ""
	if req.CheckValidity() == nil {
		t.Fail()
	}

	// url too long should fail the check too
	req.URL = token.RandomToken(cfg.MAX_URL_LEN + 1)
	if req.CheckValidity() == nil {
		t.Fail()
	}
}
