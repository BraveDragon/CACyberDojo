package charactercontroller

import (
	"CACyberDojo/model/charactermodel"
)

func GetOwnCharacters(id string) ([]charactermodel.Character, error) {
	characterIds, err := charactermodel.GetOwnCharacterIDs(id)
	if err != nil {
		return []charactermodel.Character{}, err
	}
	results := []charactermodel.Character{}
	for _, characterId := range characterIds {
		character, err := charactermodel.SearchCharacterById(characterId)
		if err != nil {
			return []charactermodel.Character{}, err
		}
		results = append(results, character)
	}

	return results, nil
}

//指定したユーザーIDに所持キャラクターを追加
func AddOwnCharacters(Userid string, results []charactermodel.Character) error {
	return charactermodel.AddOwnCharacters(Userid, results)
}
