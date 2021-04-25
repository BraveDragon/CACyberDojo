package usermodel

import "crypto/ed25519"

//User: ユーザー情報を管理.
type User struct {
	Id          string             `db:"primarykey" column:"id"`                         //ユーザーID
	Name        string             `db:"unique" column:"name" json:"name"`               //ユーザー名
	MailAddress string             `db:"unique" column:"mailAddress" json:"mailAddress"` //メールアドレス
	PassWord    string             `db:"unique" column:"passWord" json:"passWord"`       //パスワード
	PrivateKey  ed25519.PrivateKey `db:"" column:"privateKey"`                           //認証トークンの秘密鍵

}
