package charactercontroller

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/model"
	"net/http"
)

type Character struct {
	Id   int    `db:"primarykey" column:"id"`
	Name string `db:"unique" column:"name"`
}
type OwnCharacter struct {
	UserId      string `db:"" column:"userId"`
	CharacterId int    `db:"" column:"characterId"`
}

func ShowOwnCharacters_Impl(w http.ResponseWriter, r *http.Request) ([]Character, error) {
	//ユーザーを取得するためにjsonTokenを取得
	_, jsonToken, _, err := usercontroller.CheckPasetoAuth(w, r)
	if err != nil {
		return []Character{}, commonErrors.FailedToAuthorizationError()
	}

	//ログインしているユーザーを取得
	loginUser, err := usercontroller.GetOneUser(jsonToken)
	if err != nil {
		return []Character{}, commonErrors.FailedToGetUserError()
	}
	//DBに接続して所持キャラクター一覧を取得
	DB := model.Init()
	DBMap := model.NewDBMap(DB)
	OwnCharacters := []OwnCharacter{}
	DBMap.Select(&OwnCharacters, "SELECT characterId FROM owncharacters WHERE userId=?", loginUser.Id)
	Characters := []Character{}
	for _, ownCaracter := range OwnCharacters {
		Characters_tmp := []Character{}
		DBMap.Select(&Characters_tmp, "SELECT * FROM characters WHERE id=?", ownCaracter.CharacterId)
		Characters = append(Characters, Characters_tmp...)
	}

	return Characters, nil

}
