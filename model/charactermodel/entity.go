package charactermodel

//Character: キャラクターを管理.
type Character struct {
	Id       int    `db:"id, primarykey"` //キャラクターのID
	Name     string `db:"name"`           //名前
	Strength int    `db:"strength"`       //強さ
	Rarity   int    `db:"rarity"`         //レアリティ
}

//OwnCharacter: 所持キャラクターを管理.
type OwnCharacter struct {
	UserId      string `db:"userId"`      //ユーザーID
	CharacterId int    `db:"characterId"` //キャラクターのID
}
