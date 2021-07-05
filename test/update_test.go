package test

import (
	"context"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	test := &Test{}
	test.Typ = "bonjourssss"
	err := db.Update(ctx, test)
	if err != nil {
		t.Errorf("errors : %v", err)
	}
}
