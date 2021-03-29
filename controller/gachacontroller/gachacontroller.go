package gachacontroller

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"CACyberDojo/controller/usercontroller"
	"CACyberDojo/model"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Gacha struct {
	GachaId     int     `db:"" column:"gachaId"`
	CharacterId int     `db:"" column:"characterId"`
	DropRate    float64 `db:"" column:"dropRate"`
}

type GachaRequest struct {
	GachaId   int `json:"gachaId"`
	DrawTimes int `json:"drawTimes"`
}

//idに合うガチャをdrawTimes回引く
func drawGacha(id int, drawTimes int) ([]charactercontroller.Character, error) {
	DB := model.Init()
	DBMap := model.NewDBMap(DB)
	var gachaContents []Gacha
	DBMap.Select(&gachaContents, "SELECT content FROM gachas WHERE id=?", id)
	if drawTimes == 0 {
		return []charactercontroller.Character{}, commonErrors.TrytoDrawZeroTimes()
	}

	results := []charactercontroller.Character{}
	for i := 0; i < (drawTimes - 1); i++ {
		rand.Seed(time.Now().UnixNano())
		//0以上1未満の乱数を生成(結果となる)
		lottery := rand.Float64()

		for _, gachaContent := range gachaContents {
			//lotteryからcontent.DropRateの値を引いていき、lotteryが0以下になった時のcontentを結果とする
			lottery -= gachaContent.DropRate
			if lottery <= 0 {
				result := charactercontroller.Character{}
				DBMap.SelectOne(&result, "SELECT * FROM characters WHERE id=?", gachaContent.CharacterId)
				results = append(results, result)
			}

		}

	}
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
	//結果をDBに格納するためにDB,DBMapを取得
	DB := model.Init()
	DBMap := model.NewDBMap(DB)
	dbhandler, err := DBMap.Begin()
	if err != nil {
		return err
	}
	for _, result := range results {
		dbhandler.Insert(charactercontroller.OwnCharacter{UserId: loginUser.Id, CharacterId: result.Id})
	}
	dbhandler.Commit()
	return nil

}
