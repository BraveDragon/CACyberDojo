package main

import (
	"database/sql"
	"log"
	"net/http"

	"crypto/ed25519"
	"encoding/hex"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/o1egl/paseto"
)

var decoder = schema.NewDecoder()

//DB に接続
//DB : データベース本体
var DB, _ = sql.Open("mysql", "root:@APIDB")
var DBMap = &gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

//User : ユーザー情報を管理
type User struct {
	id    string `db:"ID, primarykey` //ユーザーID
	name  string `db:"name"`          //ユーザー名
	token string `db:"token"`         //認証トークン

}

func main() {
	routeCreater := mux.NewRouter()

	routeCreater.Host("http://localhost:8080")
	routeCreater.PathPrefix("http")
	routeCreater.Methods("GET", "POST", "PUT")
	routeCreater.Headers("X-Requested-With", "XMLHttpRequest")

	//エンドポイントを用意
	routeCreater.HandleFunc("/user/create/{name}", userCreate).Methods("POST").Queries("name", "{name}")
	routeCreater.HandleFunc("/user/get", userGet).Methods("GET")
	routeCreater.HandleFunc("/user/update", userUpdate).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", routeCreater))

}

func userCreate(w http.ResponseWriter, r *http.Request) {
	//パスパラメーターから新規ユーザー名を取得
	value := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	//ユーザー名
	newUser := User{}
	newUser.name = value["name"]
	//IDはUUIDで生成
	UUID, _ := uuid.NewUUID()
	newUser.id = UUID.String()
	now := time.Now()
	//ここから認証トークン生成部
	//認証トークンの生成方法は以下のサイトを参考にしている
	//URL: https://qiita.com/GpAraki/items/801cb4654ce109d49ec9
	//ユーザーIDから秘密鍵生成用のシードを生成
	b, _ := hex.DecodeString(newUser.id)
	privateKey := ed25519.PrivateKey(b)
	jsonToken := paseto.JSONToken{
		Audience:   "Audience",                       // 利用ユーザー判別するユニーク値
		Issuer:     "Issuer",                         // 利用システム
		Subject:    "WebAPI",                         // 利用機能
		Jti:        "UUID",                           // UUID
		Expiration: time.Now().Add(30 * time.Minute), // 失効日時
		IssuedAt:   now,                              // 発行日時
		NotBefore:  now,                              // 有効化日時
	}

	jsonToken.Set("KEY", "VALUE")
	footer := "FOOTER"
	tokenCreator := paseto.NewV2()
	//トークンを生成
	token, _ := tokenCreator.Sign(privateKey, jsonToken, footer)
	newUser.token = token
	dbHandler, _ := DBMap.Begin()
	//DBに追加＋反映
	dbHandler.Insert(newUser)
	dbHandler.Commit()

}

func userGet(w http.ResponseWriter, r *http.Request) {

}

func userUpdate(w http.ResponseWriter, r *http.Request) {

}
