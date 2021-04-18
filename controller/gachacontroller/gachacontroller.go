package gachacontroller

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/model/charactermodel"
	"CACyberDojo/model/gachamodel"
	"math/rand"
	"time"
)

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
func DrawGacha(id int, drawTimes int) ([]charactermodel.Character, error) {

	var gachaContents []gachamodel.Gacha
	gachamodel.SelectGacha(&gachaContents, id)
	if drawTimes == 0 {
		return []charactermodel.Character{}, commonErrors.TrytoDrawZeroTimes()
	}

	results := draw(drawTimes, gachaContents)

	return results, nil
}
