package characterhandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/handler/userhandler"
	"CACyberDojo/model/charactermodel"
	"encoding/json"
	"log"
	"net/http"
)

//ShowOwnCharacters : 所持キャラクター一覧表示のハンドラ。実際の処理はShowOwnCharactersImpl()で行う.
func ShowOwnCharacters(w http.ResponseWriter, r *http.Request) {
	Characters, err := ShowOwnCharactersImpl(w, r)

	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		_, err := w.Write([]byte("failed to get your own characters"))
		log.Print(err.Error())
		return
	}
	type result struct {
		Characters []charactermodel.Character `json:"characters"`
	}
	rawResult := result{Characters: Characters}
	resResult, err := json.Marshal(rawResult)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resResult)
	if err != nil {
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		return
	}

}

//ShowOwnCharactersImpl : ShowOwnCharactersの処理の本体.
func ShowOwnCharactersImpl(w http.ResponseWriter, r *http.Request) ([]charactermodel.Character, error) {
	//ログインしているユーザーを取得
	loginUser, err := userhandler.UserSignIn(w, r)
	if err != nil {
		return []charactermodel.Character{}, commonErrors.FailedToAuthorizationError()
	}
	//所持キャラクター一覧を取得
	Characters, err := charactercontroller.GetOwnCharacters(loginUser.Id)
	if err != nil {
		return []charactermodel.Character{}, commonErrors.FailedToSearchError()
	}

	return Characters, nil

}
