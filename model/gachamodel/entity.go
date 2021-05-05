package gachamodel

//Gacha: ガチャの中身.
type Gacha struct {
	GachaId     int     `db:"gachaId"`     //ガチャのID
	CharacterId int     `db:"characterId"` //キャラクターのID
	DropRate    float64 `db:"dropRate"`    //排出率
}
