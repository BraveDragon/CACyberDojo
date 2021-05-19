package usermodel

//User: ユーザー情報を管理.
type User struct {
	Id          string `db:"ID, primarykey"`                 //ユーザーID
	Name        string `db:"name" json:"name"`               //ユーザー名
	MailAddress string `db:"mailAddress" json:"mailAddress"` //メールアドレス
	PassWord    string `db:"password" json:"passWord"`       //パスワード
	Token       string `db:"token" json:"x-token"`           //トークン
	Score       int    `db:"score"`                          //ユーザーのスコア
}
