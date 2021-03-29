package characterhandler

import (
	"CACyberDojo/commonErrors"
	"CACyberDojo/controller/charactercontroller"
	"fmt"
	"net/http"
)

//所持キャラクター一覧表示のハンドラ
func ShowOwnCharacters(w http.ResponseWriter, r *http.Request) {
	Characters, err := charactercontroller.ShowOwnCharacters_Impl(w, r)

	if err.Error() != commonErrors.FailedToAuthorizationError().Error() {
		w.Write([]byte(fmt.Sprintf("Permission error.")))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, character := range Characters {
		w.Write([]byte(fmt.Sprintf(character.Name)))
	}

}
