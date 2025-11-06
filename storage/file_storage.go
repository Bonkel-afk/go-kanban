package storage

import (
	"encoding/json"
	"os"

	"com.bonkelbansi/go-kanban/internals/models"
)

type FileStorage struct {
	FilePath string
}

func (fs *FileStorage) LoadTasks() ([]models.Task, error) {
	f, err := os.Open(fs.FilePath)
	if os.IsNotExist(err) {
		return []models.Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tasks []models.Task
	if err := json.NewDecoder(f).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (fs *FileStorage) SaveTasks(tasks []models.Task) error {
	tmp := fs.FilePath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(tasks); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, fs.FilePath)
}
