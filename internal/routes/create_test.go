package routes

import (
	"testing"
	"time"
)

func TestRequestValidityCheck(t *testing.T) {
	// current time - 10 second
	req := request{
		ExpireAt: time.Now().Add(-time.Second * 10),
	}

	// the req should not be valid for expireAt is before current time
	if req.CheckValidity() != nil {
		t.Log("success")
	} else {
		t.Log("fail")
	}

	// current time + 10 second
	req.ExpireAt = time.Now().Add(time.Second * 10)
	if req.CheckValidity() == nil {
		t.Log("success")
	} else {
		t.Log("fail")
	}
}
