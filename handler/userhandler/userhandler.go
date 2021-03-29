package userhandler

import (
	"fmt"
	"net/http"
	"time"

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

func UserSignIn(w http.ResponseWriter, r *http.Request) {
	token, expiration, err := usercontroller.UserSignIn_Impl(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

func Refresh(w http.ResponseWriter, r *http.Request) {
	// トークンの検証(有効かどうか)
	_, jsonToken, _, err := usercontroller.CheckPasetoAuth(w, r)
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
		UserSignIn(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	err := usercontroller.UserUpdate_Impl(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
