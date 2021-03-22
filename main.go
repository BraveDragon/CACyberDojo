package main

//ユーザー関連の処理を行う
//処理については以下のサイトを参考にしている
//https://www.sohamkamani.com/golang/jwt-authentication/
//https://github.com/sohamkamani/jwt-go-example

import (
	"log"
	"net/http"

	"CACyberDojo/DataBase/gachahandler"
	"CACyberDojo/DataBase/userhandler"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func main() {
	routeCreator := mux.NewRouter()

	routeCreator.Host("https://localhost:8080")
	routeCreator.PathPrefix("https")
	routeCreator.Methods("GET", "POST", "PUT")
	routeCreator.Headers("X-Requested-With", "XMLHttpRequest")

	//エンドポイントを用意
	//ユーザー作成
	routeCreator.HandleFunc("/user/create/{name}/{mailAddress}/{passWord}",
		userhandler.UserCreate).Methods("POST").Queries("name", "mailAddress", "passWord",
		"{name}", "{mailAddress}", "{passWord}")
	//ユーザーサインイン
	routeCreator.HandleFunc("/user/signIn/{mailAddress}/{passWord}",
		userhandler.UserSignIn).Methods("GET").Queries("mailAddress", "passWord", "{mailAddress}", "{passWord}")
	//ユーザー情報取得
	routeCreator.HandleFunc("/user/get", userhandler.UserGet(userhandler.UserGet_impl)).Methods("GET")
	//トークンのリフレッシュ
	routeCreator.HandleFunc("/user/refresh", userhandler.Refresh).Methods("GET")

	//ガチャを引く
	routeCreator.HandleFunc("/gacha/draw/{gachaId}/{drawTimes}",
		gachahandler.GachaDrawHandler).Methods("GET").Queries("drawTimes", "{drawTimes}", "gachaId", "{gachaId}")

	//所持キャラクターの一覧を表示
	routeCreator.HandleFunc("/character/list", gachahandler.ShowOwnCharacters).Methods("GET")

	//ユーザー情報更新
	routeCreator.HandleFunc("/user/update", userhandler.UserUpdate).Methods("PUT").Queries("name", "{name}")
	log.Fatal(http.ListenAndServe(":8080", routeCreator))

}
