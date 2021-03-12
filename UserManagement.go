package main

//ユーザー関連の処理を行う

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
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
	id         string             `db:"ID, primarykey` //ユーザーID
	name       string             `db:"name"`          //ユーザー名
	privateKey ed25519.PrivateKey `db:"privateKey"`    //認証トークン

}

func main() {
	routeCreater := mux.NewRouter()

	routeCreater.Host("https://localhost:8080")
	routeCreater.PathPrefix("https")
	routeCreater.Methods("GET", "POST", "PUT")
	routeCreater.Headers("X-Requested-With", "XMLHttpRequest")

	//エンドポイントを用意
	//ユーザー作成
	routeCreater.HandleFunc("/user/create/{name}", userCreate).Methods("POST").Queries("name", "{name}")
	//ユーザーサインイン
	routeCreater.HandleFunc("/user/signIn", userSignIn).Methods("GET")
	//ユーザー情報取得
	routeCreater.HandleFunc("/user/get", userGet).Methods("GET")
	//ユーザー情報更新
	routeCreater.HandleFunc("/user/update", userUpdate).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", routeCreater))

}

func userCreate(w http.ResponseWriter, r *http.Request) {
	//パスパラメーターから新規ユーザー名を取得
	value := mux.Vars(r)
	//新規ユーザー
	newUser := User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbHandler, _ := DBMap.Begin()
	//TODO: DBからユーザー名を検索(既に同じのがあったらエラーを吐く)

	w.WriteHeader(http.StatusOK)

	//ユーザー名
	newUser.name = value["name"]

	//IDはUUIDで生成
	UUID, _ := uuid.NewUUID()
	newUser.id = UUID.String()

	//ここから認証トークン生成部
	//認証トークンの生成方法は以下のサイトを参考にしている
	//URL: https://qiita.com/GpAraki/items/801cb4654ce109d49ec9
	//ユーザーIDから秘密鍵生成用のシードを生成
	b, _ := hex.DecodeString(newUser.id)
	privateKey := ed25519.PrivateKey(b)
	//リフレッシュのことを考えて秘密鍵をDBに保存
	newUser.privateKey = privateKey

	//DBに追加＋反映
	dbHandler.Insert(newUser)
	dbHandler.Commit()
	//全て終わればメッセージを出して終了
	w.Write([]byte(fmt.Sprintf("User %s created", newUser.name)))

}

func userSignIn(w http.ResponseWriter, r *http.Request) {
	//ユーザー
	user := User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now()
	expiration := time.Now().Add(30 * time.Minute)
	jsonToken := paseto.JSONToken{
		Audience:   "Audience", // 利用ユーザー判別するユニーク値
		Issuer:     "Issuer",   // 利用システム
		Subject:    "WebAPI",   // 利用機能
		Jti:        "UUID",     // UUID
		Expiration: expiration, // 失効日時
		IssuedAt:   now,        // 発行日時
		NotBefore:  now,        // 有効化日時
	}

	jsonToken.Set("KEY", "VALUE")
	footer := "FOOTER"
	tokenCreator := paseto.NewV2()
	//トークンを生成
	token, err := tokenCreator.Sign(user.privateKey, jsonToken, footer)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expiration,
	})

	//全て終わればユーザー情報表示画面へ
	userGet(w, r)
}
func userGet(w http.ResponseWriter, r *http.Request) {

}

func userUpdate(w http.ResponseWriter, r *http.Request) {

}
