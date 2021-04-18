package characterhandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/model/charactermodel"
	"fmt"
	"net/http"
)

//所持キャラクター一覧表示のハンドラ
func ShowOwnCharacters(w http.ResponseWriter, r *http.Request) {
	Characters, err := ShowOwnCharacters_Impl(w, r)

	if err.Error() != commonErrors.FailedToAuthorizationError().Error() {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, character := range Characters {
		w.Write([]byte(fmt.Sprintf(character.Name)))
	}

}

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
	Characters, err := charactercontroller.GetOwnCharacters(loginUser.Id)
	if err != nil {
		return []charactermodel.Character{}, commonErrors.FailedToSearchError()
	}

	return Characters, nil

}
