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
func UserAuthorization(user *User, mailAddress string, password string) error {
	DBMap := model.NewDBMap(model.DB)
	var DBusers []User
	_, err := DBMap.Select(&DBusers, "SELECT * FROM users")
	if err != nil {
		return err
	}
	for _, DBUser := range DBusers {
		errPass := bcrypt.CompareHashAndPassword([]byte(DBUser.PassWord), []byte(password))
		errAddress := bcrypt.CompareHashAndPassword([]byte(DBUser.MailAddress), []byte(mailAddress))
		if errPass == nil && errAddress == nil {
			//両方とも一致するものがDB内にあればそれをuserに詰めて返す。返り値はnilとする
			user = &DBUser
			return nil
		}
	}
	//見つからない場合はエラーを返す
	return commonErrors.FailedToSearchError()

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
