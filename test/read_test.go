package test

import (
	"context"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	tests := []*Test{}
	err := db.Read(ctx, &tests, nil)
	if err != nil {
		t.Errorf("errors : %v", err)
	}
}
