package db

import (
	"awesomeProject2/setting"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var Link *sql.DB

func Connect(opt *setting.Setting) error {
	var e error
	Link, e = sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		opt.DbHost,
		opt.DbPort,
		opt.DbUser,
		opt.DbPass,
		opt.DbName,
	))
	if e != nil {
		return e
	}

	e = Link.Ping()
	if e != nil {
		return e
	}

	e = prepare()
	if e != nil {
		return e
	}

	e = prepareProduct()
	if e != nil {
		return e
	}
	prepareMenu()

	return nil
}
