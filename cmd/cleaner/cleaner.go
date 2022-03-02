package main

import (
	"context"
	"log"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var cfg = config.GetConfig()

func cleanExpiredRecords(records *mongo.Collection) (int64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := records.DeleteMany(ctx, bson.M{
		"expireAt": bson.M{
			"$lte": time.Now(),
		},
	})

	return res.DeletedCount, err
}

func main() {
	logger := log.Default()
	logger.SetPrefix("[ExpiredCleaner] ")

	client, err := models.Connect()
	if err != nil {
		panic(err)
	}

	records := client.Database(cfg.DATABASE).Collection(cfg.RECORD_COLLECTION)

	for {
		deletedCnt, err := cleanExpiredRecords(records)
		if err != nil {
			log.Fatal(err)
		}

		if deletedCnt != 0 {
			logger.Printf("Deleted %d expired records", deletedCnt)
		}

		<-time.After(time.Second * 10)
	}
}
