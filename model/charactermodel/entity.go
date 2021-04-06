package charactermodel

//Character: キャラクターを管理。
type Character struct {
	Id   int    `db:"primarykey" column:"id"`
	Name string `db:"unique" column:"name"`
}

//OwnCharacter: 所持キャラクターを管理。
type OwnCharacter struct {
	UserId      string `db:"" column:"userId"`
	CharacterId int    `db:"" column:"characterId"`
}
