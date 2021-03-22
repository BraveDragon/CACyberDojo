package DataBase

import (
	"database/sql"
	"log"

	"github.com/go-gorp/gorp"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() *sql.DB {
	var err error
	DB, err = sql.Open("mysql", "root:@APIDB")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return DB

}

func NewDBMap() *gorp.DbMap {
	return &gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
}
