package commonErrors

//よく使うエラーをまとめて関数化
import "errors"

func NoAuthorizationheaderError() error {
	return errors.New("There is no Authorization header!")

}

func IncorrectJsonBodyError() error {
	return errors.New("The json body is incorrect!")

}
func IncorrectTokenError() error {
	return errors.New("The token is incorrect!")

}

func InvalidSettingOfDrawerError() error {
	return errors.New("You can set only one drawer")

}

func FailedToAuthorizationError() error {
	return errors.New("Failed to authorize")
}

func FailedToCreateTokenError() error {
	return errors.New("Failed to create token")

}

func FailedToSearchError() error {
	return errors.New("Failed to search")

}

func FailedToGetUserError() error {
	return errors.New("Failed to Get a user")

}

func TrytoDrawZeroTimes() error {
	return errors.New("You try to draw gacha 0 times.")

}
