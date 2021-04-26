package handlerutil

import (
	"CACyberDojo/commonErrors"
	"encoding/json"
	"log"
	"net/http"
)

//ParseJsonBody : JSONボディから必要なデータを取得.
func ParseJsonBody(r *http.Request, decordtarget interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&decordtarget)
	if err != nil {
		return commonErrors.IncorrectJsonBodyError()
	} else {
		return nil
	}
}

//ErrorLoggingAndWriteHeader : errのnilチェック＋Log吐き＋httpステータスをw.WriteHeader()する.
func ErrorLoggingAndWriteHeader(w http.ResponseWriter, err error, httpStatus int) {
	log.Print(err.Error())
	w.WriteHeader(httpStatus)

}
