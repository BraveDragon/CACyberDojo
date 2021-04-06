package gachamodel

type Gacha struct {
	GachaId     int     `db:"" column:"gachaId"`
	CharacterId int     `db:"" column:"characterId"`
	DropRate    float64 `db:"" column:"dropRate"`
}
