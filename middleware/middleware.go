package middleware

import (
	"CACyberDojo/controller/usercontroller"
	"net/http"
	"time"
)

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			usercontroller.UserSignIn(w, r)
			next.ServeHTTP(w, r)
		})
}

//トークンのリフレッシュ用のミドルウェア
func RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// トークンの検証(有効かどうか)
			_, jsonToken, _, err := usercontroller.CheckPasetoAuth(w, r)
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
				usercontroller.UserSignIn(w, r)
				w.WriteHeader(http.StatusOK)

			}
			next.ServeHTTP(w, r)
		})

}
