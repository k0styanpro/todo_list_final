package api

import (
	"encoding/json"
	"net/http"
)

const dateLayout = "20060102"

// writeJSON кодирует data в JSON и ставит заголовок.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(data)
}
