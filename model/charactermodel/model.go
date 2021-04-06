package charactermodel

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/model"
	"CACyberDojo/model/usermodel"
)

//キャラクターIDからキャラクターを返す
func SearchCharacterById(characterId int) (Character, error) {
	DBMap := model.NewDBMap(model.DB)
	result := Character{}
	err := DBMap.SelectOne(&result, "SELECT * FROM characters WHERE id=?", characterId)
	if err != nil {
		return Character{}, commonErrors.FailedToSearchError()
	}
	return result, nil
}

//ユーザーIDから所有するキャラクター一覧を取得
func GetOwnCharacters(id string) ([]Character, error) {
	ownCharacters := []OwnCharacter{}
	DBMap := model.NewDBMap(model.DB)
	_, err := DBMap.Select(&ownCharacters, "SELECT characterId FROM owncharacters WHERE userId=?", id)
	if err != nil {
		return []Character{}, commonErrors.FailedToSearchError()
	}

	characters := []Character{}
	for _, ownCharacter := range ownCharacters {
		Characters_tmp := []Character{}
		_, err = DBMap.Select(&Characters_tmp, "SELECT * FROM characters WHERE id=?", ownCharacter.CharacterId)
		if err != nil {
			return []Character{}, commonErrors.FailedToSearchError()
		}
		characters = append(characters, Characters_tmp...)
	}

	return characters, nil

}

//所持キャラクターを追加する
func AddOwnCharacters(loginUser usermodel.User, results []Character) error {

	DBMap := model.NewDBMap(model.DB)
	dbhandler, err := DBMap.Begin()
	if err != nil {
		return err
	}
	for _, result := range results {
		err = dbhandler.Insert(OwnCharacter{UserId: loginUser.Id, CharacterId: result.Id})
		if err != nil {
			return err
		}
	}
	err = dbhandler.Commit()
	if err != nil {
		return err
	}
	return nil
}
