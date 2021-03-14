package main

//ユーザー関連の処理を行う
//処理については以下のサイトを参考にしている
//https://www.sohamkamani.com/golang/jwt-authentication/
//https://github.com/sohamkamani/jwt-go-example

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"time"

	"CACyberDojo/commonErrors"

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

//トークン生成用の定数類
//フッター
const footer = "FOOTER"

//トークンの有効期限
const expirationTime = 30 * time.Minute

//User : ユーザー情報を管理
type User struct {
	id         string             `db:"ID, primarykey` //ユーザーID
	name       string             `db:"name"`          //ユーザー名
	privateKey ed25519.PrivateKey `db:"privateKey"`    //認証トークンの秘密鍵

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
	routeCreater.HandleFunc("/user/get", userGet(userGet_impl)).Methods("GET")
	//トークンのリフレッシュ
	routeCreater.HandleFunc("/user/refresh", refresh).Methods("GET")

	//ユーザー情報更新
	routeCreater.HandleFunc("/user/update", userUpdate).Methods("PUT").Queries("name", "{name}")
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
	expiration := time.Now().Add(expirationTime)
	jsonToken := paseto.JSONToken{
		Audience:   "Audience", // 利用ユーザー判別するユニーク値
		Issuer:     "Issuer",   // 利用システム
		Subject:    "WebAPI",   // 利用機能
		Jti:        "UUID",     // UUID
		Expiration: expiration, // 失効日時
		IssuedAt:   now,        // 発行日時
		NotBefore:  now,        // 有効化日時
	}

	jsonToken.Set("ID", user.id)

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

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	(*w).Header().Add("Access-Control-Allow-Headers", "*")
}

func contains(s []string, e []string) bool {
	for _, a := range s {
		for _, b := range e {
			if a == b {
				return true
			}
		}
	}
	return false
}

//トークンの検証
func CheckPasetoAuth(w http.ResponseWriter, r *http.Request) (string, paseto.JSONToken, string, error) {
	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		//  Authorizationヘッダーがない時はエラーを返す
		w.WriteHeader(http.StatusUnauthorized)
		return "", paseto.JSONToken{}, "", commonErrors.NoAuthorizationheaderError()
	}
	tokenStr := bearerToken[7:]
	var newJsonToken paseto.JSONToken
	var newFooter string
	publicKey, _, _ := ed25519.GenerateKey(nil)
	err := paseto.NewV2().Verify(tokenStr, publicKey, &newJsonToken, &newFooter)
	if err != nil {
		return "", paseto.JSONToken{}, "", commonErrors.IncorrectTokenError()
	}

	return tokenStr, newJsonToken, newFooter, nil

}

//トークンのチェック
//ユーザー情報取得はuserGet_impl()に丸投げ
func userGet(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		enableCors(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _, _, err := CheckPasetoAuth(w, r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Permission error.")))
			return
		}

		handler(w, r)
	}
}

//ユーザー情報取得処理
func userGet_impl(w http.ResponseWriter, r *http.Request) {
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	}
	loginUser, err := getOneUser(jsonToken)
	w.Write([]byte(fmt.Sprintf(loginUser.id)))
	w.Write([]byte(fmt.Sprintf(loginUser.name)))

}

func refresh(w http.ResponseWriter, r *http.Request) {
	// トークンの検証(有効かどうか)
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		//トークンが無効ならエラーを返す
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	now := time.Now()
	//トークンの有効期限がまだ切れていない時は何もせずにそのまま返す
	if jsonToken.Expiration.After(now) == true {
		w.WriteHeader(http.StatusOK)
		return

	} else {
		//有効期限が切れていたらもう一度サインインしてトークンをリフレッシュ
		userSignIn(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}

}

//jsonTokenからユーザーを取得
func getOneUser(jsonToken paseto.JSONToken) (User, error) {
	id := jsonToken.Get("ID")
	loginUser := User{}
	err := DBMap.SelectOne(&loginUser, "SELECT * FROM user WHERE ID = ?", id)
	if err != nil {
		return loginUser, commonErrors.FailedToSearchError()

	}
	return loginUser, nil

}

func userUpdate(w http.ResponseWriter, r *http.Request) {
	// 誰がログインしているかをチェック
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//トークンから主キーのユーザーIDを取得
	loginUser, err := getOneUser(jsonToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	//パスパラメーターから新規ユーザー名を取得
	value := mux.Vars(r)
	loginUser.name = value["name"]
	dbHandler, _ := DBMap.Begin()
	_, err2 := dbHandler.Update(loginUser)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//修正したらDBに反映
	dbHandler.Commit()

}
