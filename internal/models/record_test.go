package models

import (
	"fmt"
	"testing"
	"time"
)

func TestSaveAndFindRecord(t *testing.T) {
	client, err := Connect()
	if err != nil {
		t.FailNow()
	}

	records := client.Database("shorturl_test").Collection("records")

	r, err := FindRecord(records, "testid")
	if err == nil {
		r.Delete(records)
	}

	r = Record{
		ID:       "testid",
		URL:      "http://localhost",
		ExpireAt: time.Now().Add(time.Second * 5),
	}

	err = r.Save(records)
	if err != nil {
		t.FailNow()
	}

	r, err = FindRecord(records, "testid")
	if err != nil {
		fmt.Println(2, err)
		t.FailNow()
	}
}
