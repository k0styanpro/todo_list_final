package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/k0styanpro/todo_list_final/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "title is required"})
		return
	}
	if err := normalizeDate(&task); err != nil {
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

// normalizeDate:
// 1) Если task.Date пустой — ставим сегодня.
// 2) Парсим task.Date.
// 3) Сравниваем только даты (без часов).
// 4) Если date < today и есть repeat — смещаем на NextDate.
// 5) Если date < today и нет repeat — ставим сегодня.
func normalizeDate(task *db.Task) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}
	tDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())

	if tDate.Before(today) {
		if task.Repeat != "" {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("invalid repeat: %w", err)
			}
			task.Date = next
		} else {
			task.Date = now.Format(dateLayout)
		}
	}
	return nil
}
