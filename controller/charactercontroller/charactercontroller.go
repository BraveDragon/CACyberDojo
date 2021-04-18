package charactercontroller

import (
	"CACyberDojo/model/charactermodel"
	"CACyberDojo/model/usermodel"
)

func GetOwnCharacters(id string) ([]charactermodel.Character, error) {
	return charactermodel.GetOwnCharacters(id)
}

func AddOwnCharacters(loginUser usermodel.User, results []charactermodel.Character) error {
	return charactermodel.AddOwnCharacters(loginUser, results)
}
