package userhandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/model/usermodel"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
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
	err = usercontroller.UpdateUser(loginUser)
	if err != nil {
		return err
	}

	return nil

}

//UserCreate : ユーザー作成する.
func UserCreate(w http.ResponseWriter, r *http.Request) {
	name, err := userCreateImpl(r)
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

//UserCreateImpl : userhandler.UserCreate()の処理の本体.ユーザー情報取得を行う.
func userCreateImpl(r *http.Request) (string, error) {
	jsonUser := usermodel.User{}
	//JSONボディから必要なデータを取得
	err := handlerutil.ParseJsonBody(r, &jsonUser)
	if err != nil {
		return "", err
	}
	//パスワードをハッシュ化して格納
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(jsonUser.PassWord), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	jsonUser.PassWord = string(hashedPassword)
	//メールアドレスをハッシュ化して格納
	hashedMailAddress, err := bcrypt.GenerateFromPassword([]byte(jsonUser.MailAddress), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	jsonUser.MailAddress = string(hashedMailAddress)

	//IDはUUIDで生成
	UUID, _ := uuid.NewUUID()
	id := UUID.String()
	jsonUser.Id = id

	//ここから認証トークン生成部
	//認証トークンの生成方法は以下のサイトを参考にしている
	//URL: https://qiita.com/GpAraki/items/801cb4654ce109d49ec9
	//ユーザーIDから秘密鍵生成用のシードを生成
	b, _ := hex.DecodeString(id)
	privateKey := ed25519.PrivateKey(b)
	jsonUser.PrivateKey = privateKey

	err = usermodel.CreateUser(jsonUser)
	if err != nil {
		return "", err
	}

	return jsonUser.Name, nil

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
	//ユーザーID、ユーザー名、ユーザーのスコア、ランキングを出力
	_, err = w.Write([]byte(fmt.Sprintf(loginUser.Id)))
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
	_, err = w.Write([]byte(fmt.Sprintf(loginUser.Name)))
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
	_, err = w.Write([]byte(fmt.Sprint(strconv.Itoa(loginUser.Score))))
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
	rank, err := usercontroller.GetUserRank(loginUser)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
	_, err = w.Write([]byte(fmt.Sprint(strconv.Itoa(rank))))
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

//トークン生成用の定数類
//フッター
const footer = "FOOTER"

//トークンの有効期限
const expirationTime = 30 * time.Minute

//UserSignIn : ユーザーのサインイン処理を行う.
func UserSignIn(w http.ResponseWriter, r *http.Request) {
	jsonUser := usermodel.User{}
	//JSONボディから必要なデータを取得
	err := handlerutil.ParseJsonBody(r, &jsonUser)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusBadRequest)
		return
	}

	//メールアドレスとパスワードを照合＋DBにある時のみサインインを通す
	user, err := usercontroller.UserAuthorization(jsonUser.MailAddress, jsonUser.PassWord)
	if err != nil {
		//メールアドレスとパスワードの組がDBになければエラーを返す
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusUnauthorized)
		return
	}
	now := time.Now()
	expiration := time.Now().Add(expirationTime)
	jsonToken := paseto.JSONToken{
		Expiration: expiration, // 失効日時
		IssuedAt:   now,        // 発行日時
		NotBefore:  now,        // 有効化日時
	}

	jsonToken.Set("ID", user.Id)

	tokenCreator := paseto.NewV2()

	//トークンを生成
	token, err := tokenCreator.Sign(user.PrivateKey, jsonToken, footer)

	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expiration,
	})

}
