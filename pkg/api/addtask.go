package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/k0styanpro/todo_list_final/pkg/db"
)

// POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "invalid JSON: " + err.Error()})
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "db error: " + err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"id": id})
}

// Вспомогательная коррекция даты, аналогично add и update.
func prepareDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}
	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	if !t.After(now) && task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("invalid repeat rule: %w", err)
		}
		task.Date = next
	}
	if !t.After(now) && task.Repeat == "" {
		task.Date = now.Format(dateLayout)
	}
	if t.After(now) && task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			return fmt.Errorf("invalid repeat rule: %w", err)
		}
	}
	return nil
}
