package userhandler

import (
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/model/usermodel"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

//UserUpdateImpl : ユーザー情報の更新.UserUpdate()の処理の本体.
func userUpdateImpl(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Name string `json:"name"`
	}
	// 誰がログインしているかをチェック
	loginUser, err := usercontroller.UserAuthorization(r.Header.Get("x-token"))
	if err != nil {
		return err
	}

	//jsonボディから新しい名前を取得
	rawRequest := request{}
	err = handlerutil.ParseJsonBody(r, &rawRequest)
	if err != nil {
		return err
	}

	loginUser.Name = rawRequest.Name
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
	//トークンを生成してDBに保存
	token := CreateToken(jsonUser)
	jsonUser.Token = token

	err = usercontroller.CreateUser(jsonUser)
	if err != nil {
		return "", err
	}

	return token, nil

}

//CheckPasetoAuth : トークンの検証.
func CheckPasetoAuth(w http.ResponseWriter, r *http.Request) (string, paseto.JSONToken, string, error) {

	token := r.Header.Get("x-token")
	//TODO:トークンがDBにあるかチェック
	var newJsonToken paseto.JSONToken
	var newFooter string

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
	w.Header().Set("Content-Type", "application/json")
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
	token := r.Header.Get("x-token")
	//トークンを照合＋DBにある時のみサインインを通す
	user, err := usercontroller.UserAuthorization(token)
	if err != nil {
		//トークンがDBになければエラーを返す
		return usermodel.User{}, err
	}
	return user, nil

}

//トークン生成用の定数類
//トークン生成元
const seedLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

//トークンの桁数
const tokenChrSize = 10

//CreateToken : トークンを生成する.
func CreateToken(user usermodel.User) string {
	token := make([]byte, tokenChrSize)
	//関数実行ごとにシードを変更
	rand.Seed(time.Now().UnixNano())
	for i := range token {
		token[i] = seedLetters[rand.Intn(len(seedLetters))]
	}
	return string(token)
}
