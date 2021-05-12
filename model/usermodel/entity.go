package usermodel

import "crypto/ed25519"

//User: ユーザー情報を管理.
type User struct {
	Id          string             `db:"ID, primarykey"`                 //ユーザーID
	Name        string             `db:"name" json:"name"`               //ユーザー名
	MailAddress string             `db:"mailAddress" json:"mailAddress"` //メールアドレス
	PassWord    string             `db:"password" json:"passWord"`       //パスワード
	PrivateKey  ed25519.PrivateKey `db:"privateKey"`                     //認証トークンの秘密鍵
	Score       int                `db:"score"`                          //ユーザーのスコア
}
