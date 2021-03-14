package commonErrors

//よく使うエラーをまとめて関数化
import "errors"

func NoAuthorizationheaderError() error {
	return errors.New("There is no Authorization header!")

}

func IncorrectTokenError() error {
	return errors.New("The token is incorrect!")

}

func FailedToSearchError() error {
	return errors.New("Failed to search")

}
