package api

import (
	"net/http"

	"github.com/k0styanpro/todo_list_final/pkg/db"
)

// DELETE /api/task?id=<id>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{})
}
