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
	//sql.Open()の第2引数は環境に合わせて修正すること
	DB, err = sql.Open("mysql", "MineDragon:@/cacyberdojo")
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}

func NewDBMap(DB *sql.DB) *gorp.DbMap {
	return &gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
}
