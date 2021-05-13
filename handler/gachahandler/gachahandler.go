package gachahandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/controller/gachacontroller"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/model/charactermodel"
	"encoding/json"
	"log"
	"net/http"
)

//GachaRequest : ガチャを引く時のリクエストの中身.
type GachaRequest struct {
	GachaId   int    `json:"gachaId"`
	DrawTimes int    `json:"times"`
	UserId    string `json:"id"`
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
	//何のガチャを何回引くかをリクエストで受け取る
	gachaRequest := GachaRequest{}
	err := handlerutil.ParseJsonBody(r, &gachaRequest)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		log.Print("err 1")
		return commonErrors.FailedToCreateTokenError()
	}
	results, err := gachacontroller.DrawGacha(gachaRequest.GachaId, gachaRequest.DrawTimes)
	if err != nil {
		log.Print("err 2")
		return err
	}

	err = charactercontroller.AddOwnCharacters(gachaRequest.UserId, results)
	if err != nil {
		log.Print("err 3")
		return err
	}
	loginUser, err := usercontroller.GetOneUser(gachaRequest.UserId)
	if err != nil {
		log.Print("err 4")
		return err
	}

	//ユーザーのスコアを加算
	for _, result := range results {
		err := usercontroller.AddUserScore(loginUser, result.Strength)
		if err != nil {
			log.Print("err 5")
			return err
		}
	}
	type result struct {
		Results []charactermodel.Character `json:"results"`
	}
	resResult := result{Results: results}
	res, err := json.Marshal(resResult)
	if err != nil {
		log.Print("err 6")
		return err
	}
	_, err = w.Write(res)
	if err != nil {
		log.Print("err 7")
		return err
	}
	return nil

}
