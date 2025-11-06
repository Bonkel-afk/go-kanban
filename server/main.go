package main

import (
	"log"
	"net/http"
	"os"

	"com.bonkelbansi/go-kanban/internals/web"
	"com.bonkelbansi/go-kanban/storage"
)

func main() {
	var store storage.Storage

	mode := os.Getenv("KANBAN_STORAGE")
	if mode == "mongo" {
		mongoURI := os.Getenv("KANBAN_MONGO_URI")
		ms, err := storage.NewMongoStorage(mongoURI, "kanban", "tasks")
		if err != nil {
			log.Fatalf("MongoDB-Verbindung fehlgeschlagen: %v", err)
		}
		store = ms
		log.Println("📦 MongoDB aktiv")
	} else {
		store = &storage.FileStorage{FilePath: "tasks.json"}
		log.Println("💾 Lokale JSON-Datei aktiv")
	}

	mux := web.SetupRouter(store)
	log.Println("Server läuft auf http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", mux))
}
