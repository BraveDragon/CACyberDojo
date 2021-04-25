package characterhandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/userhandler"
	"CACyberDojo/model/charactermodel"
	"fmt"
	"log"
	"net/http"
)

//ShowOwnCharacters : 所持キャラクター一覧表示のハンドラ。実際の処理はShowOwnCharactersImpl()で行う.
func ShowOwnCharacters(w http.ResponseWriter, r *http.Request) {
	Characters, err := ShowOwnCharactersImpl(w, r)

	if err != nil {
		_, err := w.Write([]byte("Failed to get your own characters."))
		log.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, character := range Characters {
		_, err := w.Write([]byte(fmt.Sprintf(character.Name)))
		if err != nil {
			log.Print(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

}

//ShowOwnCharactersImpl : ShowOwnCharactersの処理の本体.
func ShowOwnCharactersImpl(w http.ResponseWriter, r *http.Request) ([]charactermodel.Character, error) {
	//ユーザーを取得するためにjsonTokenを取得
	_, jsonToken, _, err := userhandler.CheckPasetoAuth(w, r)
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
