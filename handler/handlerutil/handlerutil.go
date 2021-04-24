package handlerutil

import (
	"CACyberDojo/commonErrors"
	"encoding/json"
	"net/http"
)

//JSONボディから必要なデータを取得
func ParseJsonBody(r *http.Request, decordtarget interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&decordtarget)
	if err != nil {
		return commonErrors.IncorrectJsonBodyError()
	} else {
		return nil
	}
}
