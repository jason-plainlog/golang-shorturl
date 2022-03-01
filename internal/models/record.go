package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	ID       string    `bson:"id"`
	URL      string    `bson:"url"`
	ExpireAt time.Time `bson:"expireAt"`
}

func FindRecord(records *mongo.Collection, id string) (Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result Record
	err := records.FindOne(ctx, bson.M{"id": id}).Decode(result)

	return result, err
}

// save records, return error if collision occurs
func (r *Record) Save(records *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result Record
	err := records.FindOne(ctx, bson.M{"id": r.ID}).Decode(result)

	if err == nil && result.ExpireAt.After(time.Now()) {
		// id in used and not yet expired
		return fmt.Errorf("id in use")
	}

	if err == mongo.ErrNoDocuments || result.ExpireAt.Before(time.Now()) {
		// id not used or expired, upsert new record
		res, err := records.UpdateOne(ctx, bson.M{"id": r.ID}, r, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}

		if res.UpsertedCount != 1 {
			return fmt.Errorf("upsert failed")
		}

		return nil
	}

	return err
}
