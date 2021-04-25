package model

import (
	"database/sql"

	"github.com/go-gorp/gorp"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

//Init : DBの初期化を行う.
func Init() error {
	var err error
	//sql.Open()の第2引数は環境に合わせて修正すること
	DB, err = sql.Open("mysql", "MineDragon:@/cacyberdojo")
	if err != nil {
		return err
	}
	return nil

}

//NewDBMap : DBMapのインスタンスを生成.
func NewDBMap(DB *sql.DB) *gorp.DbMap {
	return &gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
}
