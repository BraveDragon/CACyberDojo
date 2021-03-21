package DataBase

import (
	"database/sql"
	"log"
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
