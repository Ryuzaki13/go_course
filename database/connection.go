package database

import (
	"awesomeProject2/setting"
	"awesomeProject2/utils"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var Link *sql.DB

func Connect(opt *setting.Setting) {
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
		panic(e)
	}

	e = Link.Ping()
	if e != nil {
		panic(e)
	}

	errorList := make([]string, 0)

	errorList = append(errorList, prepareUser()...)
	errorList = append(errorList, prepareProduct()...)

	if len(errorList) > 0 {
		for _, msg := range errorList {
			utils.Logger.Println("ERROR:", msg)
		}

		panic("prepare statement failure")
	}

	LoadSession(sessionMap)
}
