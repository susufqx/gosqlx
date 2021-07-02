package test

import "time"

type Test struct {
	Id         int    `others:"pKey"`
	Typ        string `db:"type"`
	Name       string
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

func (t *Test) GetTableName() string {
	return "test"
}
