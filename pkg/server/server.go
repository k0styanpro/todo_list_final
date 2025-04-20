package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/k0styanpro/todo_list_final/pkg/api"
)

const defaultPort = 7540

// Run запускает HTTP‑сервер:
// – раздаёт статику из webDir
// – регистрирует API через api.Init()
// – слушает порт из TODO_PORT или defaultPort
func Run(webDir string) error {
	port := defaultPort
	if s := os.Getenv("TODO_PORT"); s != "" {
		if p, err := strconv.Atoi(s); err == nil && p > 0 {
			port = p
		}
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init()

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}
