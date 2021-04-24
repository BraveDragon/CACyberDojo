package userhandler

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"

	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/model/usermodel"

	"github.com/o1egl/paseto"
)

func UserUpdate_Impl(w http.ResponseWriter, r *http.Request) error {
	// 誰がログインしているかをチェック
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		return commonErrors.FailedToAuthorizationError()
	}
	//トークンから主キーのユーザーIDを取得
	loginUser, err := usercontroller.GetOneUser(jsonToken)
	if err != nil {
		return err

	}
	jsonUser := usermodel.User{}
	//jsonボディからメールアドレスとパスワードを取得
	err = json.NewDecoder(r.Body).Decode(&jsonUser)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		return commonErrors.IncorrectJsonBodyError()
	}
	loginUser.Name = jsonUser.Name
	err = usermodel.UpdateUser(loginUser)
	if err != nil {
		return err
	}

	return nil

}

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
func UserGet(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _, _, err := CheckPasetoAuth(w, r)
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
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	}
	loginUser, err := usercontroller.GetOneUser(jsonToken)
	w.Write([]byte(fmt.Sprintf(loginUser.Id)))
	w.Write([]byte(fmt.Sprintf(loginUser.Name)))

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	err := UserUpdate_Impl(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
