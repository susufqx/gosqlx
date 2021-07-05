package test

import (
	"context"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	test := &Test{Id: 0, Typ: "bonjourss"}
	err := db.Delete(ctx, test)
	if err != nil {
		t.Errorf("errors : %v", err)
	}
}

func TestDeleteKV(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	test := &Test{}
	err := db.Delete(ctx, test, "name", "nn")
	if err != nil {
		t.Errorf("errors : %v", err)
	}
}
