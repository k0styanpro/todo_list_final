package api

import (
	"encoding/json"
	"github.com/k0styanpro/todo_list_final/pkg/db"
	"net/http"
)

// PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	// Проверяем, что передан идентификатор
	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	// Заголовок обязателен
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "title is required"})
		return
	}

	// Нормализуем дату так же, как при добавлении
	if err := normalizeDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Выполняем обновление в БД
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": "db error: " + err.Error()})
		return
	}

	// Успешный ответ — пустой JSON объекта
	writeJSON(w, map[string]string{})
}
