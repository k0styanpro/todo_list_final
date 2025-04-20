// pkg/api/common.go
package api

import (
	"encoding/json"
	"net/http"
)

// dateLayout — формат даты YYYYMMDD для пакета api.
const dateLayout = "20060102"

// writeJSON устанавливает JSON‑заголовок и сериализует data в ответ.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// На случай ошибки кодирования — отправим HTTP 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
