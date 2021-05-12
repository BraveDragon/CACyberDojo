package usermodel

import (
	"CACyberDojo/model"
)

//CreateUser : ユーザーを新規作成してDBに追加.
func CreateUser(user User) error {
	dbMap := model.NewDBMap(model.DB)
	//DBと構造体を結びつける
	dbMap.AddTableWithName(User{}, "users")
	dbHandler, err := dbMap.Begin()
	if err != nil {
		return err
	}
	//DBに追加＋反映
	err = dbHandler.Insert(&user)
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
	dbMap := model.NewDBMap(model.DB)
	//DBと構造体を結びつける
	dbMap.AddTableWithName(User{}, "users")
	return dbMap.SelectOne(&user, "SELECT * FROM user WHERE ID = ?", id)
}

//UserAuthorization : ユーザーのメールアドレスとパスワードがあるかチェック.
func UserAuthorization(mailAddress string, password string) (User, error) {
	dbMap := model.NewDBMap(model.DB)
	//DBと構造体を結びつける
	dbMap.AddTableWithName(User{}, "users")
	var dbUser User
	err := dbMap.SelectOne(&dbUser, "SELECT * FROM users WHERE mailAddress=? AND password=?", mailAddress, password)
	if err != nil {
		return User{}, err
	}
	//TODO: パスワード・メールアドレスの暗号化
	// for _, DBUser := range DBusers {
	// 	errPass := bcrypt.CompareHashAndPassword([]byte(DBUser.PassWord), []byte(password))
	// 	errAddress := bcrypt.CompareHashAndPassword([]byte(DBUser.MailAddress), []byte(mailAddress))
	// 	if errPass == nil && errAddress == nil {
	// 		//両方とも一致するものがDB内にあればそれを返す
	// 		return DBUser, nil
	// 	}
	// }
	//見つからない場合はエラーを返す
	//return User{}, commonErrors.FailedToSearchError()
	return dbUser, nil
}

//UpdateUser : ユーザー名を引数の内容に更新
func UpdateUser(user User) error {
	dbMap := model.NewDBMap(model.DB)
	//DBと構造体を結びつける
	dbMap.AddTableWithName(User{}, "users")
	dbHandler, err := dbMap.Begin()
	if err != nil {
		return err
	}
	_, err = dbHandler.Update(&user)

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
	dbMap := model.NewDBMap(model.DB)
	//DBと構造体を結びつける
	dbMap.AddTableWithName(User{}, "users")
	dbHandler, err := dbMap.Begin()
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

//AddUserScore : ユーザーのスコアを加算.
func AddUserScore(user User, addScore int) error {
	dbMap := model.NewDBMap(model.DB)
	//DBと構造体を結びつける
	dbMap.AddTableWithName(User{}, "users")
	dbhandler, err := dbMap.Begin()
	if err != nil {
		return err
	}
	//ユーザーにスコアを加点
	//加点されるスコアはキャラクターの強さとなる
	user.Score += addScore
	_, err = dbhandler.Update(&user)
	if err != nil {
		return err
	}

	err = dbhandler.Commit()
	if err != nil {
		return nil
	}

	return nil

}
