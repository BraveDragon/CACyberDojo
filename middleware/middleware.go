package middleware

import (
	"CACyberDojo/handler/handlerutil"
	"CACyberDojo/handler/userhandler"
	"net/http"
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

//EnableCorsMiddleware : CORS対応用のミドルウェア.
func EnableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		//プリフライトリクエストの場合の処理
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)

	})

}
