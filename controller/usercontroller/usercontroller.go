package usercontroller

import (
	"CACyberDojo/model"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"CACyberDojo/commonErrors"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

//User: ユーザー情報を管理
type User struct {
	Id          string             `db:"primarykey" column:"id"`                         //ユーザーID
	Name        string             `db:"unique" column:"name" json:"name"`               //ユーザー名
	MailAddress string             `db:"unique" column:"mailAddress" json:"mailAddress"` //メールアドレス
	PassWord    string             `db:"unique" column:"passWord" json:"passWord"`       //パスワード
	PrivateKey  ed25519.PrivateKey `db:"" column:"privateKey"`                           //認証トークンの秘密鍵

}

//トークン生成用の定数類
//フッター
const footer = "FOOTER"

//トークンの有効期限
const expirationTime = 30 * time.Minute

func UserCreate_Impl(r *http.Request) (string, error) {
	jsonUser := User{}
	//JSONボディから必要なデータを取得
	err := json.NewDecoder(r.Body).Decode(&jsonUser)
	if err != nil {
		return "", commonErrors.IncorrectJsonBodyError()
	}
	DB := model.Init()
	DBMap := model.NewDBMap(DB)
	dbHandler, _ := DBMap.Begin()

	//IDはUUIDで生成
	UUID, _ := uuid.NewUUID()
	id := UUID.String()
	jsonUser.Id = id

	//ここから認証トークン生成部
	//認証トークンの生成方法は以下のサイトを参考にしている
	//URL: https://qiita.com/GpAraki/items/801cb4654ce109d49ec9
	//ユーザーIDから秘密鍵生成用のシードを生成
	b, _ := hex.DecodeString(id)
	privateKey := ed25519.PrivateKey(b)
	jsonUser.PrivateKey = privateKey
	//DBに追加＋反映
	dbHandler.Insert()
	dbHandler.Commit()

	return jsonUser.Name, nil

}

func UserSignIn(w http.ResponseWriter, r *http.Request) {
	jsonUser := User{}
	//JSONボディから必要なデータを取得
	err := json.NewDecoder(r.Body).Decode(&jsonUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//メールアドレスとパスワードを照合＋DBにある時のみサインインを通す
	DB := model.Init()
	DBMap := model.NewDBMap(DB)
	//ユーザー
	user := User{}
	err = DBMap.SelectOne(&user, "SELECT * FROM users WHERE mailAddress=? AND passWord=?", jsonUser.MailAddress, jsonUser.PassWord)
	if err != nil {
		//メールアドレスとパスワードの組がDBになければエラーを返す
		w.WriteHeader(http.StatusBadRequest)
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

	jsonToken.Set("ID", user.Id)

	tokenCreator := paseto.NewV2()

	//トークンを生成
	token, err := tokenCreator.Sign(user.PrivateKey, jsonToken, footer)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expiration,
	})

}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			UserSignIn(w, r)
			next.ServeHTTP(w, r)
		})
}

//jsonTokenからユーザーを取得
func GetOneUser(jsonToken paseto.JSONToken) (User, error) {
	id := jsonToken.Get("ID")
	loginUser := User{}
	DB := model.Init()
	DBMap := model.NewDBMap(DB)
	err := DBMap.SelectOne(&loginUser, "SELECT * FROM user WHERE ID = ?", id)
	if err != nil {
		return loginUser, commonErrors.FailedToSearchError()

	}
	return loginUser, nil

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

func UserUpdate_Impl(w http.ResponseWriter, r *http.Request) error {
	// 誰がログインしているかをチェック
	_, jsonToken, _, err := CheckPasetoAuth(w, r)
	if err != nil {
		return commonErrors.FailedToAuthorizationError()
	}
	//トークンから主キーのユーザーIDを取得
	loginUser, err := GetOneUser(jsonToken)
	if err != nil {
		return err

	}
	jsonUser := User{}
	//jsonボディからメールアドレスとパスワードを取得
	err = json.NewDecoder(r.Body).Decode(&jsonUser)
	if err != nil {
		//bodyの構造がおかしい時はエラーを返す
		return commonErrors.IncorrectJsonBodyError()
	}
	loginUser.Name = jsonUser.Name
	DB := model.Init()
	DBMap := model.NewDBMap(DB)
	dbHandler, _ := DBMap.Begin()
	_, err2 := dbHandler.Update(loginUser)
	if err2 != nil {
		return err2
	}

	//修正したらDBに反映
	dbHandler.Commit()
	return nil

}

func RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// トークンの検証(有効かどうか)
			_, jsonToken, _, err := CheckPasetoAuth(w, r)
			if err != nil {
				//トークンが無効ならエラーを返す
				w.WriteHeader(http.StatusUnauthorized)

			}
			now := time.Now()
			//トークンの有効期限がまだ切れていない時は何もせずにそのまま返す
			if jsonToken.Expiration.After(now) == true {
				w.WriteHeader(http.StatusOK)

			} else {
				//有効期限が切れていたらもう一度サインインしてトークンをリフレッシュ
				UserSignIn(w, r)
				w.WriteHeader(http.StatusOK)

			}
			next.ServeHTTP(w, r)
		})

}
