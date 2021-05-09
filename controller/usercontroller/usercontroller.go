package usercontroller

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/model/usermodel"
	"crypto/ed25519"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

//UserCreateImpl : userhandler.UserCreate()の処理の本体.ユーザー情報取得を行う.
func UserCreateImpl(r *http.Request) (string, error) {
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

//GetOneUser : jsonTokenからユーザーを取得.
func GetOneUser(jsonToken paseto.JSONToken) (usermodel.User, error) {
	id := jsonToken.Get("ID")
	loginUser := usermodel.User{}
	err := usermodel.GetOneUser(&loginUser, id)
	if err != nil {
		return loginUser, commonErrors.FailedToSearchError()

	}
	return loginUser, nil

}

//GetUserRank : ユーザーのランキングを取得. 処理はusermodel.GetUserRank()に丸投げ.
func GetUserRank(user usermodel.User) (int, error) {
	return usermodel.GetUserRank(user)
}

//UpdateUser : ユーザー名を引数の内容に更新. 処理はusermodel.UpdateUser()に丸投げ.
func UpdateUser(user usermodel.User) error {
	return usermodel.UpdateUser(user)
}

//AddUserScore : ユーザーのスコアを加算. 処理はusermodel.AddUserScore()に丸投げ.
func AddUserScore(user usermodel.User, addScore int) error {
	return usermodel.AddUserScore(user, addScore)
}

//UserAuthorization : ユーザー認証を行う.処理はusermodel.UserAuthorization()に丸投げ.
func UserAuthorization(mailAddress string, password string) (usermodel.User, error) {
	return usermodel.UserAuthorization(mailAddress, password)
}
