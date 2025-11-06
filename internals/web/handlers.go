package web

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"

	"com.bonkelbansi/go-kanban/internals/models"
	"com.bonkelbansi/go-kanban/storage"
)

var (
	mu    sync.Mutex
	tmpl  = template.Must(template.ParseFiles("internals/web/templates/board.html"))
	tasks []models.Task
)

func BoardHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	mu.Lock()
	var todo, doing, done []models.Task
	for _, t := range tasks {
		switch t.Status {
		case models.StatusTodo:
			todo = append(todo, t)
		case models.StatusDoing:
			doing = append(doing, t)
		case models.StatusDone:
			done = append(done, t)
		}
	}
	data := struct {
		Todo  []models.Task
		Doing []models.Task
		Done  []models.Task
	}{todo, doing, done}
	mu.Unlock()

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mu.Lock()
	id := len(tasks) + 1
	newTask := models.Task{ID: id, Title: title, Status: models.StatusTodo}
	tasks = append(tasks, newTask)
	store.SaveTasks(tasks)
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func MoveTaskHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	newStatus := models.Status(r.FormValue("status"))

	mu.Lock()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = newStatus
			break
		}
	}
	store.SaveTasks(tasks)
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
