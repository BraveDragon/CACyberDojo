package gachahandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/controller/gachacontroller"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/handler/userhandler"
	"CACyberDojo/model/charactermodel"
	"encoding/json"
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
	err := gachaDrawHandlerImpl(w, r)
	if err != nil {
		_, err = w.Write([]byte("Failed to draw gacha."))
		//エラーが出たらエラーをlogに吐く
		handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusInternalServerError)
		if err != nil {
			log.Print(err.Error())
		}
		return
	}
}

//GachaDrawHandlerImpl : GachaDrawHandler()の処理の本体.
func gachaDrawHandlerImpl(w http.ResponseWriter, r *http.Request) error {
	//ユーザーを取得するためにjsonTokenを取得
	id, _, _, err := userhandler.CheckJsonBody(r)
	if err != nil {
		return commonErrors.FailedToAuthorizationError()
	}
	//ログインしているユーザーを取得
	loginUser, err := usercontroller.GetOneUser(id)
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

	//ユーザーのスコアを加算
	for _, result := range results {
		err := usercontroller.AddUserScore(loginUser, result.Strength)
		if err != nil {
			return err
		}
	}
	type result struct {
		Results []charactermodel.Character `json:"results"`
	}
	resResult := result{Results: results}
	res, err := json.Marshal(resResult)
	if err != nil {
		return err
	}
	w.Write(res)
	return nil

}
