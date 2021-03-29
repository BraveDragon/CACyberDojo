package gachahandler

import (
	"CACyberDojo/DataBase"
	"CACyberDojo/commonErrors"
	"CACyberDojo/handler/characterhandler"
	"CACyberDojo/handler/userhandler"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type OwnCharacter struct {
	UserId      string `db:"" column:"userId"`
	CharacterId int    `db:"" column:"characterId"`
}

type Gacha struct {
	GachaId     int     `db:"primarykey" column:"gachaId"`
	CharacterId int     `db:"unique" column:"characterId"`
	DropRate    float64 `db:"" column:"dropRate"`
}

type GachaRequest struct {
	GachaId   int `json:"gachaId"`
	DrawTimes int `json:"drawTimes"`
}

//idに合うガチャをdrawTimes回引く
func drawGacha(id int, drawTimes int) ([]characterhandler.Character, error) {
	DB := DataBase.Init()
	DBMap := DataBase.NewDBMap(DB)
	var gachaContents []Gacha
	DBMap.Select(&gachaContents, "SELECT content FROM gachas WHERE id=?", id)
	if drawTimes == 0 {
		return []characterhandler.Character{}, commonErrors.TrytoDrawZeroTimes()
	}

	results := []characterhandler.Character{}
	for i := 0; i < (drawTimes - 1); i++ {
		rand.Seed(time.Now().UnixNano())
		//0以上1未満の乱数を生成(結果となる)
		lottery := rand.Float64()

		for _, gachaContent := range gachaContents {
			//lotteryからcontent.DropRateの値を引いていき、lotteryが0以下になった時のcontentを結果とする
			lottery -= gachaContent.DropRate
			if lottery <= 0 {
				result := characterhandler.Character{}
				DBMap.SelectOne(&result, "SELECT * FROM characters WHERE id=?", gachaContent.CharacterId)
				results = append(results, result)
			}

		}

	}
	return results, nil
}

//ガチャ処理のハンドラ
func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	gachaRequest := GachaRequest{}
	err := json.NewDecoder(r.Body).Decode(&gachaRequest)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := drawGacha(gachaRequest.GachaId, gachaRequest.DrawTimes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//ユーザーを取得するためにjsonTokenを取得
	_, jsonToken, _, err := userhandler.CheckPasetoAuth(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	}
	//ログインしているユーザーを取得
	loginUser, err := userhandler.GetOneUser(jsonToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//結果をDBに格納するためにDB,DBMapを取得
	DB := DataBase.Init()
	DBMap := DataBase.NewDBMap(DB)
	dbhandler, err := DBMap.Begin()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, result := range results {
		dbhandler.Insert(OwnCharacter{UserId: loginUser.Id, CharacterId: result.Id})
	}
	dbhandler.Commit()

}

//所持キャラクター一覧表示のハンドラ
func ShowOwnCharacters(w http.ResponseWriter, r *http.Request) {
	//ユーザーを取得するためにjsonTokenを取得
	_, jsonToken, _, err := userhandler.CheckPasetoAuth(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	}
	//ログインしているユーザーを取得
	loginUser, err := userhandler.GetOneUser(jsonToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//DBに接続して所持キャラクター一覧を取得
	DB := DataBase.Init()
	DBMap := DataBase.NewDBMap(DB)
	OwnCharacters := []OwnCharacter{}
	DBMap.Select(&OwnCharacters, "SELECT characterId FROM owncharacters WHERE userId=?", loginUser.Id)
	Characters := []characterhandler.Character{}
	for _, ownCaracter := range OwnCharacters {
		Characters_tmp := []characterhandler.Character{}
		DBMap.Select(&Characters_tmp, "SELECT * FROM characters WHERE id=?", ownCaracter.CharacterId)
		Characters = append(Characters, Characters_tmp...)
	}

	for _, character := range Characters {
		w.Write([]byte(fmt.Sprintf(character.Name)))
	}

}
