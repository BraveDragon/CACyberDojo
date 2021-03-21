package userhandler

import "crypto/ed25519"

type User struct {
	id          string             `db:"primarykey" column:"id"`      //ユーザーID
	name        string             `db:"unique" column:"name"`        //ユーザー名
	mailAddress string             `db:"unique" column:"mailAddress"` //メールアドレス
	passWord    string             `db:"unique" column:"passWord"`    //パスワード
	privateKey  ed25519.PrivateKey `db:"unique" column:"privateKey"`  //認証トークンの秘密鍵

}

func NewUser_Params(params User) User {
	return User{id: params.id, name: params.name, mailAddress: params.mailAddress, passWord: params.passWord, privateKey: params.privateKey}

}

func NewUser(id string, name string, mailAddress string, passWord string, privateKey ed25519.PrivateKey) User {
	return User{id: id, name: name, mailAddress: mailAddress, passWord: passWord, privateKey: privateKey}

}

func (instance User) GetID() string {
	return instance.id
}

func (instance User) GetName() string {
	return instance.name
}

func (instance User) GetPrivateKey() ed25519.PrivateKey {
	return instance.privateKey
}

func (instance User) SetName(newName string) {
	instance.name = newName
}
