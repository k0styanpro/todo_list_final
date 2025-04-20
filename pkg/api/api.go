package api

import "net/http"

func Init() {
	// без auth
	http.HandleFunc("/api/signin", signInHandler)

	// все остальные — под мидлварью auth
	http.HandleFunc("/api/nextdate", auth(nextDateHandler))
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}
