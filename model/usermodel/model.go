package usermodel

import "CACyberDojo/model"

func GetOneUser(user *User, id string) error {
	DBMap := model.NewDBMap(model.DB)
	return DBMap.SelectOne(&user, "SELECT * FROM user WHERE ID = ?", id)

}
