package main

//ユーザー関連の処理を行う
//処理については以下のサイトを参考にしている
//https://www.sohamkamani.com/golang/jwt-authentication/
//https://github.com/sohamkamani/jwt-go-example

import (
	"log"
	"net/http"

	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/characterhandler"
	"CACyberDojo/handler/gachahandler"
	"CACyberDojo/handler/userhandler"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func main() {
	//ユーザー認証をする処理用のルーター
	AuthorizationRouteCreator := mux.NewRouter()
	//ユーザー認証をしない処理用のルーター
	OtherRouteCreator := mux.NewRouter()
	//ユーザー認証とトークンのリフレッシュはミドルウェアで行う
	AuthorizationRouteCreator.Use(usercontroller.AuthorizationMiddleware)
	AuthorizationRouteCreator.Use(usercontroller.RefreshMiddleware)

	AuthorizationRouteCreator.Host("https://localhost:8080")
	AuthorizationRouteCreator.PathPrefix("https")
	AuthorizationRouteCreator.Methods("GET", "POST", "PUT")
	AuthorizationRouteCreator.Headers("X-Requested-With", "XMLHttpRequest")

	OtherRouteCreator.Host("https://localhost:8080")
	OtherRouteCreator.PathPrefix("https")
	OtherRouteCreator.Methods("GET", "POST", "PUT")
	OtherRouteCreator.Headers("X-Requested-With", "XMLHttpRequest")

	//エンドポイントを用意
	//ユーザー作成
	OtherRouteCreator.HandleFunc("/user/create", userhandler.UserCreate).Methods("POST")

	//ユーザー情報取得
	AuthorizationRouteCreator.HandleFunc("/user/get", userhandler.UserGet(userhandler.UserGet_impl)).Methods("GET")

	//ガチャを引く
	AuthorizationRouteCreator.HandleFunc("/gacha/draw", gachahandler.GachaDrawHandler).Methods("POST")

	//所持キャラクターの一覧を表示
	AuthorizationRouteCreator.HandleFunc("/character/list", characterhandler.ShowOwnCharacters).Methods("GET")

	//ユーザー情報更新
	AuthorizationRouteCreator.HandleFunc("/user/update", userhandler.UserUpdate).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", AuthorizationRouteCreator))
	log.Fatal(http.ListenAndServe(":8080", OtherRouteCreator))

}
