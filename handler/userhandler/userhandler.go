package userhandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/model/usermodel"
	"crypto/ed25519"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}
	//全て終わればメッセージを出して終了
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("User %s created", name)))
	//w.Write()のエラーチェック
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}

}

//CheckPasetoAuth : トークンの検証.
func CheckPasetoAuth(w http.ResponseWriter, r *http.Request) (string, paseto.JSONToken, string, error) {
	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		//Authorizationヘッダーがない時はエラーを返す
		w.WriteHeader(http.StatusBadRequest)
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
			handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusForbidden)
			_, err := w.Write([]byte("permission error"))
			if err != nil {
				log.Print(err.Error())
			}

			return
		}

		handler(w, r)
	}
}

//UserGetImpl : ユーザー情報取得処理を行う.
func UserGetImpl(w http.ResponseWriter, r *http.Request) {
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusForbidden)
		_, err = w.Write([]byte("permission error"))
		if err != nil {
			log.Print(err.Error())
		}
		return
	}
	//ログインしているユーザーを取得
	loginUser, err := usercontroller.GetOneUser(jsonToken)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
	//ユーザーID、ユーザー名、ユーザーのスコアを出力
	//TODO: ユーザーのランキング機能を追加
	_, err = w.Write([]byte(fmt.Sprintf(loginUser.Id)))
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
	_, err = w.Write([]byte(fmt.Sprintf(loginUser.Name)))
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
	_, err = w.Write([]byte(fmt.Sprintf(strconv.Itoa(loginUser.Score))))
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
}

//UserUpdate : ユーザー情報の更新.処理の中身はUserUpdateImpl()に丸投げ.
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	err := UserUpdateImpl(w, r)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}

}
