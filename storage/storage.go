package storage

import (
	"log"
	"os"

	"com.bonkelbansi/go-kanban/internals/models"
)

type Storage interface {
	LoadTasks() ([]models.Task, error)
	SaveTasks([]models.Task) error
	ResetDemo([]models.Task) error
}

const DataFile = "tasks.json"

var Store Storage

// initStore wird in main() aufgerufen
func InitStore() (*MongoStore, error) {
	mode := os.Getenv("KANBAN_STORAGE") // "file" oder "mongo"

	if mode == "mongo" {
		uri := os.Getenv("KANBAN_MONGO_URI")
		if uri == "" {
			uri = "mongodb://localhost:27017"
		}
		log.Println("Nutze Mongo-Storage:", uri)
		mStore, err := NewMongoStore(uri)
		if err != nil {
			return nil, err
		}
		Store = mStore
		return mStore, nil
	}

	log.Println("Nutze File-Storage:", DataFile)
	Store = &FileStore{Path: DataFile}
	return nil, nil
}
