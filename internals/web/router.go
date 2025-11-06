package web

import (
	"net/http"

	"com.bonkelbansi/go-kanban/storage"
)

func SetupRouter(store storage.Storage) *http.ServeMux {
	mux := http.NewServeMux()

	// Board-Seite
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		BoardHandler(w, r, store)
	})

	// Task hinzufügen
	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		AddTaskHandler(w, r, store)
	})

	// Task verschieben
	mux.HandleFunc("/move", func(w http.ResponseWriter, r *http.Request) {
		MoveTaskHandler(w, r, store)
	})

	// Statische Dateien (CSS etc.)
	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("internals/web/static")),
		),
	)

	return mux
}
