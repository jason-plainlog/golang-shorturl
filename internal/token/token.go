package token

import (
	"math/rand"
	"url-shortener/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// generate random tokens
func RandomToken(l int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, l)
	rand.Read(b)
	for i := range b {
		if b[i] >= byte(len(charset)) {
			b[i] = byte(rand.Intn(len(charset)))
		}
		b[i] = charset[b[i]]
	}
	return string(b)
}

// generates unused tokens offline to boost performance, noted that the token may be used by other servers, therefore still need to check if collision occurs
func GenToken(records *mongo.Collection, tokenChan chan string) {
	for {
		token := RandomToken(6)

		_, err := models.FindRecord(records, token)
		if err == mongo.ErrNoDocuments {
			tokenChan <- token
		} else if err != nil {
			panic(err)
		}
	}
}
