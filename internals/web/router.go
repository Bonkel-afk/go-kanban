package web

import (
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", BoardHandler)
	mux.HandleFunc("/add", AddTaskHandler)
	mux.HandleFunc("/move", MoveTaskHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internals/web/static"))))
	return mux
}
