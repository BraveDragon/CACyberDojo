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

type Characters struct {
	Id   int    `db:"primarykey" column:"id"`
	Name string `db:"unique" column:"name"`
}

type OwnCharacters struct {
	UserId      string `db:"" column:"userId"`
	CharacterId int    `db:"" column:"characterId"`
}

type Content struct {
	CharacterId int     `json:"id"`
	DropRate    float64 `json:"dropRate"`
}

type Gachas struct {
	Id      int    `db:"primarykey" column:"id"`
	Content string `db:"" column:"content"`
}

//idに合うガチャをdrawTimes回引く
func drawGacha(id int, drawTimes int) ([]Characters, error) {
	DB := DataBase.Init()
	DBMap := DataBase.NewDBMap(DB)
	var gacha Gachas
	DBMap.SelectOne(&gacha, "SELECT content FROM gachas WHERE id=?", id)
	if drawTimes == 0 {
		return []Characters{}, commonErrors.TrytoDrawZeroTimes()
	}
	byteContent := []byte(gacha.Content)
	contents := []Content{}
	err := json.Unmarshal(byteContent, &contents)
	if err != nil {
		return []Characters{}, err
	}
	results := []Characters{}
	for i := 0; i < (drawTimes - 1); i++ {
		rand.Seed(time.Now().UnixNano())
		//0以上1未満の乱数を生成(結果となる)
		lottery := rand.Float64()

		for _, content := range contents {
			//lotteryからcontent.DropRateの値を引いていき、lotteryが0以下になった時のcontentを結果とする
			lottery -= content.DropRate
			if lottery <= 0 {
				result := Characters{}
				DBMap.SelectOne(&result, "SELECT * FROM characters WHERE id=?", content.CharacterId)
				results = append(results, result)
			}

		}

	}
	return results, nil
}

func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	value := mux.Vars(r)
	gachaId, _ := strconv.Atoi(value["gachaId"])
	drawTimes, _ := strconv.Atoi(value["drawTimes"])
	results, err := drawGacha(gachaId, drawTimes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, jsonToken, _, err := userhandler.CheckPasetoAuth(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	}
	//ログインしているユーザーを取得
	loginUser, err := userhandler.GetOneUser(jsonToken)
	DB := DataBase.Init()
	DBMap := DataBase.NewDBMap(DB)
	dbhandler, err := DBMap.Begin()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, result := range results {
		dbhandler.Insert(OwnCharacters{UserId: loginUser.Id, CharacterId: result.Id})
	}
	dbhandler.Commit()

}
