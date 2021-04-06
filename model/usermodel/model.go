package usermodel

import "CACyberDojo/model"

//ユーザーを新規作成してDBに追加
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

//IDからユーザーを取得
func GetOneUser(user *User, id string) error {
	DBMap := model.NewDBMap(model.DB)
	return DBMap.SelectOne(&user, "SELECT * FROM user WHERE ID = ?", id)
}

//ユーザーのメールアドレスとパスワードがあるかチェック
func UserAuthorization(user *User, mailAddress string, password string) error {
	DBMap := model.NewDBMap(model.DB)
	return DBMap.SelectOne(&user, "SELECT * FROM users WHERE mailAddress=? AND passWord=?", mailAddress, password)
}

//ユーザー名を引数の内容に更新
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
