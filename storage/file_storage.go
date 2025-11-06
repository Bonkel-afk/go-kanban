package storage

import (
	"encoding/json"
	"os"

	"go-kanban/internal/models"
)

type FileStore struct {
	Path string
}

func (f *FileStore) LoadTasks() ([]models.Task, error) {
	_, err := os.Stat(f.Path)
	if os.IsNotExist(err) {
		return []models.Task{}, nil
	}
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tasks []models.Task
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (f *FileStore) SaveTasks(ts []models.Task) error {
	tmp := f.Path + ".tmp"
	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ts); err != nil {
		return err
	}
	return os.Rename(tmp, f.Path)
}

func (f *FileStore) ResetDemo(demo []models.Task) error {
	return f.SaveTasks(demo)
}
