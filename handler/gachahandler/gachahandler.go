package gachahandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/gachacontroller"
	"fmt"
	"net/http"
)

//ガチャ処理のハンドラ
func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	err := gachacontroller.GachaDrawHandler_Impl(w, r)
	if err.Error() == commonErrors.FailedToAuthorizationError().Error() {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
