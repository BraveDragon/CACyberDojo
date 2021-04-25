package gachamodel

//Gacha: ガチャの中身.
type Gacha struct {
	GachaId     int     `db:"" column:"gachaId"`
	CharacterId int     `db:"" column:"characterId"`
	DropRate    float64 `db:"" column:"dropRate"`
}
