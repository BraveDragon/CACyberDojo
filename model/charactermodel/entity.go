package charactermodel

//Character: キャラクターを管理.
type Character struct {
	Id       uint   `db:"primarykey" column:"id"`
	Name     string `db:"unique" column:"name"`
	Strength uint   `db:"" column:"strength"`
}

//OwnCharacter: 所持キャラクターを管理.
type OwnCharacter struct {
	UserId      string `db:"" column:"userId"`
	CharacterId int    `db:"" column:"characterId"`
}
