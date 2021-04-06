package charactermodel

type Character struct {
	Id   int    `db:"primarykey" column:"id"`
	Name string `db:"unique" column:"name"`
}
type OwnCharacter struct {
	UserId      string `db:"" column:"userId"`
	CharacterId int    `db:"" column:"characterId"`
}
