package usermodel

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/model"

	"golang.org/x/crypto/bcrypt"
)

//CreateUser : ユーザーを新規作成してDBに追加.
func CreateUser(user User) error {
	DBMap := model.NewDBMap(model.DB)
	dbHandler, err := DBMap.Begin()
	if err != nil {
		return err
	}
	//DBに追加＋反映
	err = dbHandler.Insert(user)
	if err != nil {
		return err
	}
	err = dbHandler.Commit()
	if err != nil {
		return err
	}

	return nil

}

//GetOneUser : IDからユーザーを取得.
func GetOneUser(user *User, id string) error {
	DBMap := model.NewDBMap(model.DB)
	return DBMap.SelectOne(&user, "SELECT * FROM user WHERE ID = ?", id)
}

//UserAuthorization : ユーザーのメールアドレスとパスワードがあるかチェック.
func UserAuthorization(mailAddress string, password string) (User, error) {
	DBMap := model.NewDBMap(model.DB)
	var DBusers []User
	_, err := DBMap.Select(&DBusers, "SELECT * FROM users")
	if err != nil {
		return User{}, err
	}
	for _, DBUser := range DBusers {
		errPass := bcrypt.CompareHashAndPassword([]byte(DBUser.PassWord), []byte(password))
		errAddress := bcrypt.CompareHashAndPassword([]byte(DBUser.MailAddress), []byte(mailAddress))
		if errPass == nil && errAddress == nil {
			//両方とも一致するものがDB内にあればそれを返す
			return DBUser, nil
		}
	}
	//見つからない場合はエラーを返す
	return User{}, commonErrors.FailedToSearchError()

}

//UpdateUser : ユーザー名を引数の内容に更新
func UpdateUser(user User) error {
	DBMap := model.NewDBMap(model.DB)
	dbHandler, err := DBMap.Begin()
	if err != nil {
		return err
	}
	_, err = dbHandler.Update(user)

	if err != nil {
		return err
	}
	//修正したらDBに反映
	err = dbHandler.Commit()
	if err != nil {
		return err
	}
	return nil

}

//GetUserRank : ユーザーのランキングを取得.
func GetUserRank(user User) (int, error) {
	DBMap := model.NewDBMap(model.DB)
	dbHandler, err := DBMap.Begin()
	if err != nil {
		return -1, err
	}
	var allUsers []User

	_, err = dbHandler.Select(&allUsers, "SELECT * FROM users ORDER BY score DESC")
	if err != nil {
		return -1, err
	}
	var rank int
	for i, allUser := range allUsers {
		if user.Id == allUser.Id {
			rank = i
			break
		}

	}
	//for文は0からカウントするため、ランキングとして表示するために1を足す
	rank += 1
	return rank, nil
}
