package model

import (
	"database/sql"
	"log"

	"github.com/go-gorp/gorp"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init() *sql.DB {
	var err error
	//DBがnilの時のみDBを生成
	if db == nil {
		db, err = sql.Open("mysql", "root:@APIDB")
		if err != nil {
			log.Fatal(err)
			return nil
		}
	}
	return db

}

func NewDBMap(DB *sql.DB) *gorp.DbMap {
	return &gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
}
