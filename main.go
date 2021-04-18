package main

//ユーザー関連の処理を行う
//処理については以下のサイトを参考にしている
//https://www.sohamkamani.com/golang/jwt-authentication/
//https://github.com/sohamkamani/jwt-go-example

import (
	"log"
	"net/http"

	"CACyberDojo/handler/characterhandler"
	"CACyberDojo/handler/gachahandler"
	"CACyberDojo/handler/userhandler"
	"CACyberDojo/middleware"
	"CACyberDojo/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func main() {
	err := model.Init()
	if err != nil {
		log.Fatal(err)
	}
	//ユーザー認証をする処理用のルーター
	AuthorizationRouteCreator := mux.NewRouter()
	//ユーザー認証をしない処理用のルーター
	OtherRouteCreator := mux.NewRouter()
	//ユーザー認証とトークンのリフレッシュはミドルウェアで行う
	AuthorizationRouteCreator.Use(middleware.AuthorizationMiddleware)
	AuthorizationRouteCreator.Use(middleware.RefreshMiddleware)
	//CORS対応もミドルウェアで行う
	AuthorizationRouteCreator.Use(middleware.EnableCorsMiddleware)

	AuthorizationRouteCreator.Host("https://localhost:8080")
	AuthorizationRouteCreator.PathPrefix("https")
	AuthorizationRouteCreator.Headers("X-Requested-With", "XMLHttpRequest")
	//CORS対応もミドルウェアで行う
	OtherRouteCreator.Use(middleware.EnableCorsMiddleware)

	OtherRouteCreator.Host("https://localhost:8080")
	OtherRouteCreator.PathPrefix("https")
	OtherRouteCreator.Headers("X-Requested-With", "XMLHttpRequest")

	//エンドポイントを用意
	//ユーザー作成
	OtherRouteCreator.HandleFunc("/user/create", userhandler.UserCreate).Methods("POST")

	//ユーザー情報取得
	AuthorizationRouteCreator.HandleFunc("/user/get", userhandler.UserGet(userhandler.UserGet_impl)).Methods("GET")

	//ユーザー情報更新
	AuthorizationRouteCreator.HandleFunc("/user/update", userhandler.UserUpdate).Methods("PUT")

	//ガチャを引く
	AuthorizationRouteCreator.HandleFunc("/gacha/draw", gachahandler.GachaDrawHandler).Methods("POST")

	//所持キャラクターの一覧を表示
	AuthorizationRouteCreator.HandleFunc("/character/list", characterhandler.ShowOwnCharacters).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", AuthorizationRouteCreator))
	log.Fatal(http.ListenAndServe(":8080", OtherRouteCreator))

}
