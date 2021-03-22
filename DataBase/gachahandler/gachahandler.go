package gachahandler

import (
	"CACyberDojo/DataBase"
	"CACyberDojo/DataBase/userhandler"
	"CACyberDojo/commonErrors"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"strconv"

	"github.com/gorilla/mux"
)

type Character struct {
	Id   int    `db:"primarykey" column:"id"`
	Name string `db:"unique" column:"name"`
}

type OwnCharacter struct {
	UserId      string `db:"" column:"userId"`
	CharacterId int    `db:"" column:"characterId"`
}

type Content struct {
	CharacterId int     `json:"id"`
	DropRate    float64 `json:"dropRate"`
}

type Gacha struct {
	Id      int    `db:"primarykey" column:"id"`
	Content string `db:"" column:"content"`
}

//idに合うガチャをdrawTimes回引く
func drawGacha(id int, drawTimes int) ([]Character, error) {
	DB := DataBase.Init()
	DBMap := DataBase.NewDBMap(DB)
	var gacha Gacha
	DBMap.SelectOne(&gacha, "SELECT content FROM gachas WHERE id=?", id)
	if drawTimes == 0 {
		return []Character{}, commonErrors.TrytoDrawZeroTimes()
	}
	byteContent := []byte(gacha.Content)
	contents := []Content{}
	err := json.Unmarshal(byteContent, &contents)
	if err != nil {
		return []Character{}, err
	}
	results := []Character{}
	for i := 0; i < (drawTimes - 1); i++ {
		rand.Seed(time.Now().UnixNano())
		//0以上1未満の乱数を生成(結果となる)
		lottery := rand.Float64()

		for _, content := range contents {
			//lotteryからcontent.DropRateの値を引いていき、lotteryが0以下になった時のcontentを結果とする
			lottery -= content.DropRate
			if lottery <= 0 {
				result := Character{}
				DBMap.SelectOne(&result, "SELECT * FROM characters WHERE id=?", content.CharacterId)
				results = append(results, result)
			}

		}

	}
	return results, nil
}

//ガチャ処理のハンドラ
func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)
	gachaId, _ := strconv.Atoi(value["gachaId"])
	drawTimes, _ := strconv.Atoi(value["drawTimes"])
	results, err := drawGacha(gachaId, drawTimes)
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
	Characters := []Character{}
	for _, ownCaracter := range OwnCharacters {
		Characters_tmp := []Character{}
		DBMap.Select(&Characters_tmp, "SELECT * FROM characters WHERE id=?", ownCaracter.CharacterId)
		Characters = append(Characters, Characters_tmp...)
	}

	for _, character := range Characters {
		w.Write([]byte(fmt.Sprintf(character.Name)))
	}

}
