package api

import (
	"net/http"

	"github.com/k0styanpro/todo_list_final/pkg/db"
)

const defaultTasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// GET /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tasks, err := db.Tasks(defaultTasksLimit)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJSON(w, TasksResp{Tasks: tasks})
}
