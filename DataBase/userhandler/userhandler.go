package userhandler

import (
	"CACyberDojo/DataBase"
	"CACyberDojo/commonErrors"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/o1egl/paseto"
)

//User: 1ユーザー情報を管理
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

//トークン生成用の定数類
//フッター
const footer = "FOOTER"

//トークンの有効期限
const expirationTime = 30 * time.Minute

func UserCreate(w http.ResponseWriter, r *http.Request) {
	//パスパラメーターから新規ユーザー名を取得
	value := mux.Vars(r)

	err := json.NewDecoder(r.Body).Decode(&User{})
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	DBMap := DataBase.NewDBMap()
	dbHandler, _ := DBMap.Begin()

	w.WriteHeader(http.StatusOK)
	name := value["name"]
	mailAddress := value["mailAddress"]
	passWord := value["passWord"]

	//IDはUUIDで生成
	UUID, _ := uuid.NewUUID()
	id := UUID.String()

	//ここから認証トークン生成部
	//認証トークンの生成方法は以下のサイトを参考にしている
	//URL: https://qiita.com/GpAraki/items/801cb4654ce109d49ec9
	//ユーザーIDから秘密鍵生成用のシードを生成
	b, _ := hex.DecodeString(id)
	privateKey := ed25519.PrivateKey(b)

	//DBに追加＋反映
	dbHandler.Insert(NewUser(id, name, mailAddress, passWord, privateKey))
	dbHandler.Commit()
	//全て終わればメッセージを出して終了
	w.Write([]byte(fmt.Sprintf("User %s created", name)))

}

func UserSignIn(w http.ResponseWriter, r *http.Request) {
	//ユーザー
	user := User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now()
	expiration := time.Now().Add(expirationTime)
	jsonToken := paseto.JSONToken{
		Audience:   "Audience", // 利用ユーザー判別するユニーク値
		Issuer:     "Issuer",   // 利用システム
		Subject:    "WebAPI",   // 利用機能
		Jti:        "UUID",     // UUID
		Expiration: expiration, // 失効日時
		IssuedAt:   now,        // 発行日時
		NotBefore:  now,        // 有効化日時
	}

	jsonToken.Set("ID", user.GetID())

	tokenCreator := paseto.NewV2()

	//トークンを生成
	token, err := tokenCreator.Sign(user.GetPrivateKey(), jsonToken, footer)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expiration,
	})

}

//トークンの検証
func CheckPasetoAuth(w http.ResponseWriter, r *http.Request) (string, paseto.JSONToken, string, error) {
	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		//  Authorizationヘッダーがない時はエラーを返す
		w.WriteHeader(http.StatusUnauthorized)
		return "", paseto.JSONToken{}, "", commonErrors.NoAuthorizationheaderError()
	}
	tokenStr := bearerToken[7:]
	var newJsonToken paseto.JSONToken
	var newFooter string
	publicKey, _, _ := ed25519.GenerateKey(nil)
	err := paseto.NewV2().Verify(tokenStr, publicKey, &newJsonToken, &newFooter)
	if err != nil {
		return "", paseto.JSONToken{}, "", commonErrors.IncorrectTokenError()
	}

	return tokenStr, newJsonToken, newFooter, nil

}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	(*w).Header().Add("Access-Control-Allow-Headers", "*")
}

//トークンのチェック
//ユーザー情報取得はuserGet_impl()に丸投げ
func UserGet(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		enableCors(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _, _, err := CheckPasetoAuth(w, r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Permission error.")))
			return
		}

		handler(w, r)
	}
}

//ユーザー情報取得処理
func UserGet_impl(w http.ResponseWriter, r *http.Request) {
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	}
	loginUser, err := getOneUser(jsonToken)
	w.Write([]byte(fmt.Sprintf(loginUser.GetID())))
	w.Write([]byte(fmt.Sprintf(loginUser.GetName())))

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// トークンの検証(有効かどうか)
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		//トークンが無効ならエラーを返す
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	now := time.Now()
	//トークンの有効期限がまだ切れていない時は何もせずにそのまま返す
	if jsonToken.Expiration.After(now) == true {
		w.WriteHeader(http.StatusOK)
		return

	} else {
		//有効期限が切れていたらもう一度サインインしてトークンをリフレッシュ
		UserSignIn(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}

}

//jsonTokenからユーザーを取得
func getOneUser(jsonToken paseto.JSONToken) (User, error) {
	id := jsonToken.Get("ID")
	loginUser := NewUser_Params(User{})
	DBMap := DataBase.NewDBMap()
	err := DBMap.SelectOne(&loginUser, "SELECT * FROM user WHERE ID = ?", id)
	if err != nil {
		return loginUser, commonErrors.FailedToSearchError()

	}
	return loginUser, nil

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	// 誰がログインしているかをチェック
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//トークンから主キーのユーザーIDを取得
	loginUser, err := getOneUser(jsonToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	//パスパラメーターから新規ユーザー名を取得
	value := mux.Vars(r)
	loginUser.SetName(value["name"])
	DBMap := DataBase.NewDBMap()
	dbHandler, _ := DBMap.Begin()
	_, err2 := dbHandler.Update(loginUser)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//修正したらDBに反映
	dbHandler.Commit()

}
