package userhandler

import (
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/model/usermodel"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

//UserUpdateImpl : ユーザー情報の更新.UserUpdate()の処理の本体.
func userUpdateImpl(w http.ResponseWriter, r *http.Request) error {
	// 誰がログインしているかをチェック
	jsonUser := usermodel.User{}
	jsonUser.Id = r.Header.Get("id")
	//jsonボディからIDを取得
	err := handlerutil.ParseJsonBody(r, &jsonUser)
	if err != nil {
		return err
	}
	loginUser, err := usercontroller.GetOneUser(jsonUser.Id)
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
	token, err := userCreateImpl(r)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}
	//全て終わればメッセージを出して終了
	w.WriteHeader(http.StatusOK)

	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}
	//トークンをjson形式で返す
	type result struct {
		Token string `json:"token"`
	}
	rawResult := result{Token: token}
	resResult, err := json.Marshal(rawResult)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resResult)
	//w.Write()のエラーチェック
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}

}

//秘密鍵生成用のシード(128桁)
//TODO: 文字列を途中で折り返す方法を見つける
const secretKey = "276538ba123456091749759837598027127498375957987902740982774983748a276538ba12345609174975983759802712749837595798790274098273748a"

//UserCreateImpl : userhandler.UserCreate()の処理の本体.ユーザー情報取得を行う.
func userCreateImpl(r *http.Request) (string, error) {
	jsonUser := usermodel.User{}
	//JSONボディから必要なデータを取得
	err := handlerutil.ParseJsonBody(r, &jsonUser)
	if err != nil {
		return "", err
	}
	//TODO:パスワード・メールアドレスのハッシュ化
	//パスワードをハッシュ化して格納
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(jsonUser.PassWord), bcrypt.DefaultCost)
	// if err != nil {
	// 	return "", err
	// }
	// jsonUser.PassWord = string(hashedPassword)
	//メールアドレスをハッシュ化して格納
	// hashedMailAddress, err := bcrypt.GenerateFromPassword([]byte(jsonUser.MailAddress), bcrypt.DefaultCost)
	// if err != nil {
	// 	return "", err
	// }
	// jsonUser.MailAddress = string(hashedMailAddress)

	//IDはUUIDで生成
	UUID, _ := uuid.NewUUID()
	id := UUID.String()
	jsonUser.Id = id

	//ここから認証トークン生成部
	//認証トークンの生成方法は以下のサイトを参考にしている
	//URL: https://qiita.com/GpAraki/items/801cb4654ce109d49ec9
	b, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", err
	}
	privateKey := ed25519.PrivateKey(b)

	jsonUser.PrivateKey = privateKey

	err = usermodel.CreateUser(jsonUser)
	if err != nil {
		return "", err
	}
	//トークンを生成して返す
	token, err := CreateToken(jsonUser)
	if err != nil {
		return "", err
	}

	return token, nil

}

//CheckPasetoAuth : トークンの検証.
func CheckPasetoAuth(w http.ResponseWriter, r *http.Request) (string, paseto.JSONToken, string, error) {

	token := r.Header.Get("x-token")
	var newJsonToken paseto.JSONToken
	var newFooter string

	//公開鍵を生成
	// b, err := hex.DecodeString(secretKey)
	// if err != nil {
	// 	return "", paseto.JSONToken{}, "", err
	// }
	// publicKey := ed25519.PrivateKey(b).Public()
	// //TODO:トークンを検証
	// err = paseto.NewV2().Verify(token, publicKey, &newJsonToken, &newFooter)
	// if err != nil {
	// 	log.Print("error when verify")
	// 	return "", paseto.JSONToken{}, "", err
	// }

	return token, newJsonToken, newFooter, nil

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
	loginUser, err := UserSignIn(w, r)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusForbidden)
		_, err = w.Write([]byte("permission error"))
		if err != nil {
			log.Print(err.Error())
		}
		return
	}

	rank, err := usercontroller.GetUserRank(loginUser)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
	type result struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Score string `json:"score"`
		Rank  string `json:"rank"`
	}
	rawResult := result{Id: loginUser.Id, Name: loginUser.Name, Score: strconv.Itoa(loginUser.Score), Rank: strconv.Itoa(rank)}
	resResult, err := json.Marshal(rawResult)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}
	//ユーザーID、ユーザー名、ユーザーのスコア、ランキングをjson形式で出力
	_, err = w.Write(resResult)

	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
	}
}

//UserUpdate : ユーザー情報の更新.処理の中身はUserUpdateImpl()に丸投げ.
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	err := userUpdateImpl(w, r)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

//UserSignIn : ユーザーのサインイン処理を行う.
func UserSignIn(w http.ResponseWriter, r *http.Request) (usermodel.User, error) {
	loginUser := usermodel.User{}
	loginUser.MailAddress = r.Header.Get("mailaddress")
	loginUser.PassWord = r.Header.Get("password")
	loginUser.Id = r.Header.Get("id")
	//メールアドレスとパスワードを照合＋DBにある時のみサインインを通す
	user, err := usercontroller.UserAuthorization(loginUser.Id, loginUser.MailAddress, loginUser.PassWord)
	if err != nil {
		//メールアドレスとパスワードの組がDBになければエラーを返す
		return usermodel.User{}, err
	}
	return user, nil

}

//トークン生成用の定数類
//フッター
const footer = "eyJraWQiOiAiMTIzNDUifQ"

//トークンの有効期限
const expirationTime = 30 * time.Minute

//CreateToken : トークンを生成する.
func CreateToken(user usermodel.User) (string, error) {
	now := time.Now()
	expiration := time.Now().Add(expirationTime)
	jsonToken := paseto.JSONToken{
		Expiration: expiration, // 失効日時
		IssuedAt:   now,        // 発行日時
		NotBefore:  now,        // 有効化日時
	}
	jsonToken.Set("id", user.Id)

	tokenCreator := paseto.NewV2()
	//トークンを生成
	token, err := tokenCreator.Sign(user.PrivateKey, jsonToken, footer)
	return token, err
}
