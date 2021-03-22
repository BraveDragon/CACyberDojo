package main

//ユーザー関連の処理を行う
//処理については以下のサイトを参考にしている
//https://www.sohamkamani.com/golang/jwt-authentication/
//https://github.com/sohamkamani/jwt-go-example

import (
	"log"
	"net/http"

	"CACyberDojo/DataBase"
	"CACyberDojo/DataBase/userhandler"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

//DB に接続
//DB : データベース本体
var DB = DataBase.Init()
var DBMap = DataBase.NewDBMap()

func main() {
	routeCreater := mux.NewRouter()

	routeCreater.Host("https://localhost:8080")
	routeCreater.PathPrefix("https")
	routeCreater.Methods("GET", "POST", "PUT")
	routeCreater.Headers("X-Requested-With", "XMLHttpRequest")

	//エンドポイントを用意
	//ユーザー作成
	routeCreater.HandleFunc("/user/create/{name}/{mailAddress}/{passWord}",
		userhandler.UserCreate).Methods("POST").Queries("name", "mailAddress", "passWord",
		"{name}", "{mailAddress}", "{passWord}")
	//ユーザーサインイン
	routeCreater.HandleFunc("/user/signIn/{mailAddress}/{passWord}",
		userhandler.UserSignIn).Methods("GET").Queries("mailAddress", "passWord", "{mailAddress}", "{passWord}")
	//ユーザー情報取得
	routeCreater.HandleFunc("/user/get", userhandler.UserGet(userhandler.UserGet_impl)).Methods("GET")
	//トークンのリフレッシュ
	routeCreater.HandleFunc("/user/refresh", userhandler.Refresh).Methods("GET")

	//ユーザー情報更新
	routeCreater.HandleFunc("/user/update", userhandler.UserUpdate).Methods("PUT").Queries("name", "{name}")
	log.Fatal(http.ListenAndServe(":8080", routeCreater))

}
