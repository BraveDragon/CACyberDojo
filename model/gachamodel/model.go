package gachamodel

import "CACyberDojo/model"

//SelectGacha : idに合うガチャの中身を抽出.
func SelectGacha(contents *[]Gacha, id int) error {
	dbMap := model.NewDBMap(model.DB)
	//DBと構造体を結びつける
	dbMap.AddTableWithName(Gacha{}, "gachas")
	_, err := dbMap.Select(&contents, "SELECT * FROM gachas WHERE id=?", id)
	return err
}
