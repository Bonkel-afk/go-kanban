package storage

import (
	"encoding/json"
	"os"

	"com.bonkelbansi/go-kanban/internals/models"
)

func LoadTasks(filename string) ([]models.Task, error) {
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return []models.Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tasks []models.Task
	err = json.NewDecoder(f).Decode(&tasks)
	return tasks, err
}

func SaveTasks(filename string, tasks []models.Task) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(tasks)
}
