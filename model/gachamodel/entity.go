package gachamodel

//Gacha: ガチャの中身.
type Gacha struct {
	GachaId     int     `db:"" column:"gachaId"`     //ガチャのID
	CharacterId int     `db:"" column:"characterId"` //キャラクターのID
	DropRate    float64 `db:"" column:"dropRate"`    //排出率
}
