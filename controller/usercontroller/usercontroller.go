package usercontroller

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/model/usermodel"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

//トークン生成用の定数類
//フッター
const footer = "FOOTER"

//トークンの有効期限
const expirationTime = 30 * time.Minute

func UserCreate_Impl(r *http.Request) (string, error) {
	jsonUser := usermodel.User{}
	//JSONボディから必要なデータを取得
	err := json.NewDecoder(r.Body).Decode(&jsonUser)
	if err != nil {
		return "", commonErrors.IncorrectJsonBodyError()
	}

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

	usermodel.CreateUser(jsonUser)

	return jsonUser.Name, nil

}

func UserSignIn(w http.ResponseWriter, r *http.Request) {
	jsonUser := usermodel.User{}
	//JSONボディから必要なデータを取得
	err := json.NewDecoder(r.Body).Decode(&jsonUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//ユーザー
	user := usermodel.User{}
	//メールアドレスとパスワードを照合＋DBにある時のみサインインを通す
	err = usermodel.UserAuthorization(&user, jsonUser.MailAddress, jsonUser.PassWord)
	if err != nil {
		//メールアドレスとパスワードの組がDBになければエラーを返す
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusBadRequest)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expiration,
	})

}

//jsonTokenからユーザーを取得
func GetOneUser(jsonToken paseto.JSONToken) (usermodel.User, error) {
	id := jsonToken.Get("ID")
	loginUser := usermodel.User{}
	err := usermodel.GetOneUser(&loginUser, id)
	if err != nil {
		return loginUser, commonErrors.FailedToSearchError()

	}
	return loginUser, nil

}
