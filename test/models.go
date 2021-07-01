package test

import "time"

type Test struct {
	Id         int    `gosqlx:"pKey"`
	Typ        string `gosqlx:"column:type"`
	Name       string
	CreateTime time.Time `gosqlx:"column:create_time"`
	UpdateTime time.Time `gosqlx:"column:update_time"`
}

func (t *Test) GetTableName() string {
	return "test"
}
