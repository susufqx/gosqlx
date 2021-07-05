package gosqlx

import (
	"context"
	"testing"
	"time"

	"github.com/susufqx/gosqlx/test"
)

func TestGetTableName(t *testing.T) {
	expect := "test"
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	bm := []*test.Test{}
	r := getTableName(ctx, &bm)

	if r != expect {
		t.Errorf("WRONG")
	}
}
