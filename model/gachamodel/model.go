package gachamodel

import "CACyberDojo/model"

//SelectGacha : idに合うガチャの中身を抽出.
func SelectGacha(contents *[]Gacha, id int) error {
	DBMap := model.NewDBMap(model.DB)
	_, err := DBMap.Select(&contents, "SELECT * FROM gachas WHERE id=?", id)
	return err
}
