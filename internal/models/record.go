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
	res := records.FindOne(ctx, bson.M{"id": id})
	if res.Err() != nil {
		return result, res.Err()
	}

	err := res.Decode(&result)

	return result, err
}

var ErrIdInUse = fmt.Errorf("id in use")

// save records, return error if collision occurs
func (r *Record) Save(records *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result Record
	err := records.FindOne(ctx, bson.M{"id": r.ID}).Decode(result)

	if err == nil && result.ExpireAt.After(time.Now()) {
		// id in used and not yet expired
		return ErrIdInUse
	}

	if err == mongo.ErrNoDocuments || result.ExpireAt.Before(time.Now()) {
		// id not used or expired, upsert new record
		res, err := records.UpdateOne(ctx, bson.M{"id": r.ID}, bson.M{"$set": r}, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}

		if res.UpsertedCount+res.ModifiedCount != 1 {
			return fmt.Errorf("upsert failed")
		}

		return nil
	}

	return err
}

func (r *Record) Delete(records *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := records.DeleteOne(ctx, bson.M{"id": r.ID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("record not exist")
	}

	return nil
}
