package gachahandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/controller/gachacontroller"
	"CACyberDojo/controller/usercontroller"
	"encoding/json"
	"fmt"
	"net/http"
)

type GachaRequest struct {
	GachaId   int `json:"gachaId"`
	DrawTimes int `json:"drawTimes"`
}

//ガチャ処理のハンドラ
func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	err := GachaDrawHandler_Impl(w, r)
	if err.Error() == commonErrors.FailedToAuthorizationError().Error() {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func GachaDrawHandler_Impl(w http.ResponseWriter, r *http.Request) error {
	//ユーザーを取得するためにjsonTokenを取得
	_, jsonToken, _, err := usercontroller.CheckPasetoAuth(w, r)
	if err != nil {
		return commonErrors.FailedToAuthorizationError()
	}
	//ログインしているユーザーを取得
	loginUser, err := usercontroller.GetOneUser(jsonToken)
	if err != nil {
		return err
	}
	//何のガチャを何回引くかをリクエストで受け取る
	gachaRequest := GachaRequest{}
	err = json.NewDecoder(r.Body).Decode(&gachaRequest)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		return commonErrors.FailedToCreateTokenError()
	}
	results, err := gachacontroller.DrawGacha(gachaRequest.GachaId, gachaRequest.DrawTimes)
	if err != nil {
		return err
	}

	err = charactercontroller.AddOwnCharacters(loginUser.Id, results)
	if err != nil {
		return err
	}
	return nil

}
