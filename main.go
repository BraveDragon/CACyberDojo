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
)

func init() {
	log.SetFlags(log.Llongfile)

}

func main() {
	err := model.Init()
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()

	//ユーザー認証をする処理用のルーター
	authorizationRouteCreator := router.PathPrefix("").Subrouter()

	//ユーザー認証をしない処理用のルーター
	otherRouteCreator := router.PathPrefix("").Subrouter()

	//ユーザー認証はミドルウェアで行う
	authorizationRouteCreator.Use(middleware.AuthorizationMiddleware)
	//CORS対応もミドルウェアで行う
	authorizationRouteCreator.Use(middleware.EnableCorsMiddleware)

	//CORS対応もミドルウェアで行う
	otherRouteCreator.Use(middleware.EnableCorsMiddleware)

	//エンドポイントを用意
	//ユーザー作成
	otherRouteCreator.HandleFunc("/user/create", userhandler.UserCreate).Methods("POST")

	//ユーザー情報取得
	authorizationRouteCreator.HandleFunc("/user/get", userhandler.UserGet(userhandler.UserGetImpl)).Methods("GET")

	//ユーザー情報更新
	authorizationRouteCreator.HandleFunc("/user/update", userhandler.UserUpdate).Methods("PUT")

	//ガチャを引く
	authorizationRouteCreator.HandleFunc("/gacha/draw", gachahandler.GachaDrawHandler).Methods("POST")

	//所持キャラクターの一覧を表示
	authorizationRouteCreator.HandleFunc("/character/list", characterhandler.ShowOwnCharacters).Methods("GET")
	//サーバーを起動
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Fatal(srv.ListenAndServe())

}
