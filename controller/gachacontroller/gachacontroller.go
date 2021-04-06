package gachacontroller

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/model/charactermodel"
	"CACyberDojo/model/gachamodel"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type GachaRequest struct {
	GachaId   int `json:"gachaId"`
	DrawTimes int `json:"drawTimes"`
}

type Drawer func(drawTimes int, gachaContents []gachamodel.Gacha) []charactermodel.Character

//確変も何も行わない普通のガチャ。drawGachaのデフォルト。
func draw(drawTimes int, gachaContents []gachamodel.Gacha) []charactermodel.Character {
	results := []charactermodel.Character{}
	for i := 0; i < (drawTimes - 1); i++ {
		rand.Seed(time.Now().UnixNano())
		//0以上1未満の乱数を生成(結果となる)
		lottery := rand.Float64()

		for _, gachaContent := range gachaContents {
			//lotteryからcontent.DropRateの値を引いていき、lotteryが0以下になった時のcontentを結果とする
			lottery -= gachaContent.DropRate
			if lottery <= 0 {
				result, err := charactermodel.SearchCharacterById(gachaContent.CharacterId)
				if err != nil {
					return []charactermodel.Character{}
				}
				results = append(results, result)
			}

		}

	}
	return results
}

//idに合うガチャをdrawTimes回引く
func drawGacha(id int, drawTimes int) ([]charactermodel.Character, error) {

	var gachaContents []gachamodel.Gacha
	gachamodel.SelectGacha(&gachaContents, id)
	if drawTimes == 0 {
		return []charactermodel.Character{}, commonErrors.TrytoDrawZeroTimes()
	}

	results := draw(drawTimes, gachaContents)

	return results, nil
}

func GachaDrawHandler_Impl(w http.ResponseWriter, r *http.Request) error {
	gachaRequest := GachaRequest{}
	err := json.NewDecoder(r.Body).Decode(&gachaRequest)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		return commonErrors.FailedToCreateTokenError()
	}
	results, err := drawGacha(gachaRequest.GachaId, gachaRequest.DrawTimes)
	if err != nil {
		return err
	}
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
	err = charactermodel.AddOwnCharacters(loginUser, results)
	if err != nil {
		return err
	}
	return nil

}
