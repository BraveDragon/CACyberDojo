package gachahandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/controller/gachacontroller"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/handler/userhandler"
	"log"
	"net/http"
)

//GachaRequest : ガチャを引く時のリクエストの中身.
type GachaRequest struct {
	GachaId   int `json:"gachaId"`
	DrawTimes int `json:"drawTimes"`
}

//GachaDrawHandler : ガチャ処理のハンドラ.処理本体はGachaDrawHandlerImpl()に丸投げ.
func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	err := GachaDrawHandlerImpl(w, r)
	if err != nil {
		_, err = w.Write([]byte("Failed to draw gacha."))
		//エラーが出たらエラーをlogに吐く
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusBadRequest)
		if err != nil {
			log.Print(err.Error())
		}
		return
	}
}

//GachaDrawHandlerImpl : GachaDrawHandler()の処理の本体.
func GachaDrawHandlerImpl(w http.ResponseWriter, r *http.Request) error {
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
