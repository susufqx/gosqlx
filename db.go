package gosqlx

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type operation struct {
	p PreparerContext
	// options *string
	// model   BaseModelInterface
}

// DB : database struct
type DB struct {
	*operation
}

// Tx : transaction struct
type Tx struct {
	*operation
}

// Open : create the db connections in pool
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{&operation{p: db}}, nil
}

// BeginTx : begin transactions
func (d *DB) BeginTx() (*Tx, error) {
	db, ok := d.p.(*sqlx.DB)
	if !ok {
		return nil, errors.New("WRONG DB TYPE")
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	return &Tx{&operation{p: tx}}, nil
}

// Commit : commit the queries, if fail, rollback
func (t *Tx) Commit() error {
	tx, ok := t.p.(*sqlx.Tx)
	if !ok {
		return errors.New("WRONG Tx TYPE")
	}

	err := tx.Commit()
	if err != nil {
		err = tx.Rollback()
	}

	return err
}

// Read : find by the options
func (p *operation) Read(ctx context.Context, baseModels interface{}, options map[string]interface{}) error {
	return Read(ctx, p.p, baseModels, options)
}

// ReadPageSort : find by order and offset
func (p *operation) ReadPageSort(ctx context.Context, baseModels interface{}, options map[string]interface{}, size, offset int, orderKey, orderDire string) error {
	return ReadPageSort(ctx, p.p, baseModels, options, size, offset, orderKey, orderDire)
}

// Save : if the db model exists, update the content,
// or insert a new to db
func (p *operation) Save(ctx context.Context, baseModel BaseModelInterface) error {
	return Save(ctx, p.p, baseModel)
}

// Create : create a new, no judge if the model exists
func (p *operation) Create(ctx context.Context, baseModel BaseModelInterface) error {
	return Create(ctx, p.p, baseModel)
}

// Update : update the data without judging the model's existance
func (p *operation) Update(ctx context.Context, baseModel BaseModelInterface) error {
	return Update(ctx, p.p, baseModel)
}

// // UpdateMap : map is to record the update key-values
// func (p *operation) UpdateMap(ctx context.Context, tableName string, qm map[string]interface{}, cm map[string]interface{}) error {
// 	return UpdateMap(ctx, p.p, tableName, qm, cm)
// }

// Delete : delete the data by primary keys by default
func (p *operation) Delete(ctx context.Context, options ...interface{}) error {
	return Delete(ctx, p.p, options...)
}
