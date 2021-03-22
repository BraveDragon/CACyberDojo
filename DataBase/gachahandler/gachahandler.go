package gachahandler

import "CACyberDojo/DataBase"

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
func DrawGacha(id int, drawTimes int) {
	DB := DataBase.Init()
	DBMap := DataBase.NewDBMap(DB)
	var gacha Gachas
	DBMap.SelectOne(&gacha, "SELECT content FROM gachas WHERE id=?", id)

}
