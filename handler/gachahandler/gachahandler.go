package gachahandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/controller/gachacontroller"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/handler/userhandler"
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
		_, err = w.Write([]byte("Permission error."))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func GachaDrawHandler_Impl(w http.ResponseWriter, r *http.Request) error {
	//ユーザーを取得するためにjsonTokenを取得
	_, jsonToken, _, err := userhandler.CheckPasetoAuth(w, r)
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
	err = handlerutil.ParseJsonBody(r, &gachaRequest)
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
