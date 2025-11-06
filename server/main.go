package main

import (
	"log"
	"net/http"

	"com.bonkelbansi/go-kanban/internals/web"
	"com.bonkelbansi/go-kanban/storage"
)

func main() {
	mStore, err := storage.InitStore()
	if err != nil {
		log.Fatal("Storage init fehlgeschlagen:", err)
	}
	if mStore != nil {
		defer mStore.Close()
	}

	web.Load()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internals/web/static"))))
	http.HandleFunc("/", web.BoardHandler)
	http.HandleFunc("/add", web.AddHandler)
	http.HandleFunc("/move", web.MoveHandler)
	http.HandleFunc("/delete", web.DeleteHandler)
	http.HandleFunc("/reset", web.ResetHandler)

	log.Println("Server läuft auf http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
