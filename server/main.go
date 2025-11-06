package main

import (
	"log"
	"net/http"
	"os"

	"com.bonkelbansi/go-kanban/internals/web"
	"com.bonkelbansi/go-kanban/storage"
)

func main() {
	mode := os.Getenv("KANBAN_STORAGE")
	if mode == "" {
		mode = "file"
	}

	var store storage.Storage

	switch mode {
	case "mongo":
		uri := os.Getenv("KANBAN_MONGO_URI")
		if uri == "" {
			uri = "mongodb://localhost:27017"
		}
		ms, err := storage.NewMongoStorage(uri, "kanban", "tasks")
		if err != nil {
			log.Fatalf("MongoDB-Verbindung fehlgeschlagen: %v", err)
		}
		store = ms
		log.Println("📦 Storage-Modus: MongoDB (", uri, ")")

	default:
		store = &storage.FileStorage{FilePath: "tasks.json"}
		log.Println("💾 Storage-Modus: FileStorage (tasks.json)")
	}

	// 👉 ganz wichtig: Tasks aus dem aktiven Storage laden
	web.InitTasks(store)

	mux := web.SetupRouter(store)

	log.Println("Server läuft auf http://localhost:9090")
	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}
}
