package handlerutil

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

//ParseJsonBody : JSONボディから必要なデータを取得.
func ParseJsonBody(r *http.Request, decordtarget interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&decordtarget)
	if err != nil {
		return err
	} else {
		return nil
	}
}

//ErrorLoggingAndWriteHeader : errのnilチェック＋Log吐き＋httpステータスをw.WriteHeader()する.
func ErrorLoggingAndWriteHeader(w http.ResponseWriter, err error, httpStatus int) {
	_, src, l, _ := runtime.Caller(1)
	log.Printf("%s:%d %v", src, l, err)
	w.WriteHeader(httpStatus)

}
