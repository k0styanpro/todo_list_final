package api

import (
	"encoding/json"
	"net/http"

	"github.com/k0styanpro/todo_list_final/pkg/db"
)

// PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}
	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "title is required"})
		return
	}
	if err := prepareDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{})
}
