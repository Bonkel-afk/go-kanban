package storage

import "com.bonkelbansi/go-kanban/internals/models"

// Storage ist ein Interface, das von FileStorage und MongoStorage implementiert wird.
type Storage interface {
	LoadTasks() ([]models.Task, error)
	SaveTasks([]models.Task) error
}
