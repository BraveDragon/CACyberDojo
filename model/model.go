package model

import (
	"database/sql"
	"log"

	"github.com/go-gorp/gorp"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() error {
	var err error
	DB, err = sql.Open("mysql", "root:@APIDB")
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}

func NewDBMap(DB *sql.DB) *gorp.DbMap {
	return &gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
}
