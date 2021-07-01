package test

import (
	"context"
	"testing"
	"time"

	"github.com/susufqx/gosqlx"

	_ "github.com/lib/pq"
)

var db *gosqlx.DB

func init() {
	var err error
	var t *testing.T
	db, err = gosqlx.Open(conf.driverName, conf.driverDataSource)
	if err != nil {
		t.Errorf("errors : %v", err)
	}
}

func TestCreate(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	testModel := Test{}
	err := db.Create(ctx, &testModel)
	if err != nil {
		t.Errorf("errors : %v", err)
	}
}

func TestSave(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	testModelOne, testModeltwo := Test{}, Test{Id: 1, Typ: "xx", Name: "nn"}
	testModelOne.Typ = "test"
	testModelOne.Name = "hello"

	tx, err := db.BeginTx()
	err = tx.Save(ctx, &testModelOne)
	if err != nil {
		t.Errorf("errors : %v", err)
	}

	tx.Save(ctx, &testModeltwo)
	if err != nil {
		t.Errorf("errors : %v", err)
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("errors : %v", err)
	}
}