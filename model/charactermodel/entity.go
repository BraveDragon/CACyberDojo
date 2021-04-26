package charactermodel

//Character: キャラクターを管理.
type Character struct {
	Id       int    `db:"primarykey" column:"id"` //キャラクターのID
	Name     string `db:"unique" column:"name"`   //名前
	Strength int    `db:"" column:"strength"`     //強さ
	Rarity   int    `db:"" column:"rarity"`       //レアリティ
}

//OwnCharacter: 所持キャラクターを管理.
type OwnCharacter struct {
	UserId      string `db:"" column:"userId"`      //ユーザーID
	CharacterId int    `db:"" column:"characterId"` //キャラクターのID
}
