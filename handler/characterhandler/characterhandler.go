package characterhandler

type Character struct {
	Id   int    `db:"primarykey" column:"id"`
	Name string `db:"unique" column:"name"`
}
