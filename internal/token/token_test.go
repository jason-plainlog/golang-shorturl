package token

import (
	"testing"
	"url-shortener/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestGenRandomToken(t *testing.T) {
	for i := 6; i < 16; i++ {
		token := RandomToken(i)
		if len(token) != i {
			t.Fail()
		}
	}
}

func TestGenToken(t *testing.T) {
	tokenChan := make(chan string, 10)

	client, err := models.Connect()
	if err != nil {
		t.FailNow()
	}

	records := client.Database("shorturl").Collection("records")
	go GenToken(records, tokenChan)

	for i := 0; i < 100; i++ {
		token := <-tokenChan
		_, err := models.FindRecord(records, token)
		if err != mongo.ErrNoDocuments {
			t.Fail()
		}
	}
}
