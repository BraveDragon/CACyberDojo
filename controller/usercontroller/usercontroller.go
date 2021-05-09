package usercontroller

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/model/usermodel"

	"github.com/o1egl/paseto"
)

//GetOneUser : jsonTokenからユーザーを取得.
func GetOneUser(jsonToken paseto.JSONToken) (usermodel.User, error) {
	id := jsonToken.Get("id")
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
