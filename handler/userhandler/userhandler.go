package userhandler

import (
	"crypto/ed25519"
	"fmt"
	"log"
	"net/http"

	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/model/usermodel"

	"github.com/o1egl/paseto"
)

//UserUpdateImpl : ユーザー情報の更新.UserUpdate()の処理の本体.
func UserUpdateImpl(w http.ResponseWriter, r *http.Request) error {
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
	err = handlerutil.ParseJsonBody(r, &jsonUser)
	if err != nil {
		return err
	}
	loginUser.Name = jsonUser.Name
	err = usermodel.UpdateUser(loginUser)
	if err != nil {
		return err
	}

	return nil

}

//UserCreate : ユーザー作成する.
func UserCreate(w http.ResponseWriter, r *http.Request) {
	name, err := usercontroller.UserCreateImpl(r)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//全て終わればメッセージを出して終了
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("User %s created", name)))
	//w.Write()のエラーチェック
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

//CheckPasetoAuth : トークンの検証.
func CheckPasetoAuth(w http.ResponseWriter, r *http.Request) (string, paseto.JSONToken, string, error) {
	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		//Authorizationヘッダーがない時はエラーを返す
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

//UserGet : トークンのチェックを行う.ユーザー情報取得はUserGetImpl()に丸投げ.
func UserGet(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _, _, err := CheckPasetoAuth(w, r)
		if err != nil {
			_, err := w.Write([]byte("Permission error."))
			log.Printf(err.Error())
			if err != nil {
				log.Printf(err.Error())
				w.WriteHeader(http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(w, r)
	}
}

//UserGetImpl : ユーザー情報取得処理を行う.
func UserGetImpl(w http.ResponseWriter, r *http.Request) {
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		_, err = w.Write([]byte("Permission error."))
		if err != nil {
			log.Printf(err.Error())
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	loginUser, err := usercontroller.GetOneUser(jsonToken)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
	_, err = w.Write([]byte(fmt.Sprintf(loginUser.Id)))
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
	_, err = w.Write([]byte(fmt.Sprintf(loginUser.Name)))
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
}

//UserUpdate : ユーザー情報の更新.処理の中身はUserUpdateImpl()に丸投げ.
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	err := UserUpdateImpl(w, r)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
