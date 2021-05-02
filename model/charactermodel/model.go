package charactermodel

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/model"
	"CACyberDojo/model/usermodel"
)

//SearchCharacterById : キャラクターIDからキャラクターを返す.
func SearchCharacterById(characterId int) (Character, error) {
	DBMap := model.NewDBMap(model.DB)
	result := Character{}
	err := DBMap.SelectOne(&result, "SELECT * FROM characters WHERE id=?", characterId)
	if err != nil {
		return Character{}, commonErrors.FailedToSearchError()
	}
	return result, nil
}

//GetOwnCharacterIDs : ユーザーIDから所有するキャラクターのキャラクターIDを全て取得.
func GetOwnCharacterIDs(id string) ([]int, error) {
	ownCharacters := []OwnCharacter{}
	DBMap := model.NewDBMap(model.DB)
	_, err := DBMap.Select(&ownCharacters, "SELECT characterId FROM owncharacters WHERE userId=?", id)
	if err != nil {
		return []int{-1}, commonErrors.FailedToSearchError()
	}
	results := []int{}
	for _, ownCharacter := range ownCharacters {
		results = append(results, ownCharacter.CharacterId)
	}
	return results, nil

}

//AddOwnCharacters : 所持キャラクターを追加する.
func AddOwnCharacters(userid string, characters []Character) error {

	DBMap := model.NewDBMap(model.DB)
	dbhandler, err := DBMap.Begin()
	if err != nil {
		return err
	}
	//ログインしているユーザーを取得
	var loginUser usermodel.User
	err = usermodel.GetOneUser(&loginUser, userid)
	if err != nil {
		return err
	}
	for _, character := range characters {
		err := dbhandler.Insert(OwnCharacter{UserId: userid, CharacterId: character.Id})
		if err != nil {
			return err
		}
		//ユーザーにスコアを加点
		//加点されるスコアはキャラクターの強さとなる
		usermodel.AddUserScore(loginUser, character.Strength)

	}
	err = dbhandler.Commit()
	if err != nil {
		return err
	}
	return nil
}
