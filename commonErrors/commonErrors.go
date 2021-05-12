package commonErrors

//よく使うエラーをまとめて関数化
import "errors"

//NoAuthorizationheaderError : Authorizationヘッダーがない時のエラー.
func NoAuthorizationheaderError() error {
	return errors.New("there is no Authorization header")

}

//IncorrectJsonBodyError : Jsonボディが不正な時のエラー.
func IncorrectJsonBodyError() error {
	return errors.New("the json body is incorrect")

}

//IncorrectTokenError : トークンが不正な時のエラー.
func IncorrectTokenError() error {
	return errors.New("the token is incorrect")

}

//FailedToAuthorizationError : ユーザー認証に失敗した時のエラー.
func FailedToAuthorizationError() error {
	return errors.New("failed to authorize")
}

//FailedToCreateTokenError : トークンの生成に失敗した時のエラー.
func FailedToCreateTokenError() error {
	return errors.New("failed to create token")

}

//FailedToSearchError : DBでの検索が失敗した時のエラー.
func FailedToSearchError() error {
	return errors.New("failed to search")

}

//FailedToGetUserError : ユーザー取得に失敗した時のエラー.
func FailedToGetUserError() error {
	return errors.New("failed to Get a user")

}

//TrytoDrawZeroTimes : 0回ガチャを引こうとしたときのエラー.
func TrytoDrawZeroTimes() error {
	return errors.New("you try to draw gacha 0 times.")

}
