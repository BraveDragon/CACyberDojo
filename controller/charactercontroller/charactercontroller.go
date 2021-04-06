package charactercontroller

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/model/charactermodel"
	"net/http"
)

func ShowOwnCharacters_Impl(w http.ResponseWriter, r *http.Request) ([]charactermodel.Character, error) {
	//ユーザーを取得するためにjsonTokenを取得
	_, jsonToken, _, err := usercontroller.CheckPasetoAuth(w, r)
	if err != nil {
		return []charactermodel.Character{}, commonErrors.FailedToAuthorizationError()
	}

	//ログインしているユーザーを取得
	loginUser, err := usercontroller.GetOneUser(jsonToken)
	if err != nil {
		return []charactermodel.Character{}, commonErrors.FailedToGetUserError()
	}
	//所持キャラクター一覧を取得
	Characters, err := charactermodel.GetOwnCharacters(loginUser.Id)
	if err != nil {
		return []charactermodel.Character{}, commonErrors.FailedToSearchError()
	}

	return Characters, nil

}
