package gosqlx

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type BaseModelInterface interface {
	GetTableName() string
}

type PreparerContext interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}
