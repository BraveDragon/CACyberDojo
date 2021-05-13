package middleware

import (
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/handler/userhandler"
	"net/http"
	"time"
)

//AuthorizationMiddleware : ユーザー認証用のミドルウェア.
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, err := userhandler.UserSignIn(w, r)
			if err != nil {
				//サインインに失敗すればエラーをログに記録
				handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
}

//RefreshMiddleware : トークンのリフレッシュ用のミドルウェア.
func RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// トークンの検証(有効かどうか)
			_, jsonToken, _, err := userhandler.CheckPasetoAuth(w, r)
			if err != nil {
				//トークンが無効ならエラーを返す
				handlerutil.ErrorLoggingAndWriteHeader(w, err, http.StatusUnauthorized)
				return
			}
			now := time.Now()
			//トークンの有効期限がまだ切れていない時は何もせずにそのまま返す
			if jsonToken.Expiration.After(now) {
				//何もしない
			} else {
				//有効期限が切れていたらもう一度サインインしてトークンをリフレッシュ
				userhandler.UserSignIn(w, r)

			}
			next.ServeHTTP(w, r)
		})

}

//EnableCorsMiddleware : CORS対応用のミドルウェア.
func EnableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		//プリフライトリクエストの場合の処理
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)

	})

}
