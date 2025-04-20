package main

import (
	"log"
	"os"

	"github.com/k0styanpro/todo_list_final/pkg/db"
	"github.com/k0styanpro/todo_list_final/pkg/server"
)

func main() {
	// Инициализируем БД (используя TODO_DBFILE или по умолчанию scheduler.db)
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	// Запускаем HTTP‑сервер (статические файлы + API)
	webDir := "web"
	if err := server.Run(webDir); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
