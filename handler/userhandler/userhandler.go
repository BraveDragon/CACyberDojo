package userhandler

import (
	"fmt"
	"net/http"

	"CACyberDojo/controller/usercontroller"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	name, err := usercontroller.UserCreate_Impl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//全て終わればメッセージを出して終了
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("User %s created", name)))

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	(*w).Header().Add("Access-Control-Allow-Headers", "*")
}

//トークンのチェック
//ユーザー情報取得はuserGet_impl()に丸投げ
func UserGet(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		enableCors(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _, _, err := usercontroller.CheckPasetoAuth(w, r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Permission error.")))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r)
	}
}

//ユーザー情報取得処理
func UserGet_impl(w http.ResponseWriter, r *http.Request) {
	_, jsonToken, _, err := usercontroller.CheckPasetoAuth(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	}
	loginUser, err := usercontroller.GetOneUser(jsonToken)
	w.Write([]byte(fmt.Sprintf(loginUser.Id)))
	w.Write([]byte(fmt.Sprintf(loginUser.Name)))

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	err := usercontroller.UserUpdate_Impl(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
