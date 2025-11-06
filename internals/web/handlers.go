package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"com.bonkelbansi/go-kanban/internals/models"
	"com.bonkelbansi/go-kanban/storage"
)

type BoardData struct {
	Todo  []models.Task
	Doing []models.Task
	Done  []models.Task
}

var (
	tmpl   = template.Must(template.ParseFiles("internal/web/templates/board.html"))
	tasks  []models.Task
	mu     sync.Mutex
	nextID = 1
)

func next() int {
	id := nextID
	nextID++
	return id
}

func Load() {
	mu.Lock()
	defer mu.Unlock()

	loaded, err := storage.Store.LoadTasks()
	if err != nil {
		log.Println("Fehler beim Laden:", err)
	}
	tasks = loaded
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
	if len(tasks) == 0 {
		tasks = []models.Task{
			{ID: next(), Title: "Go installieren", Status: models.StatusTodo},
			{ID: next(), Title: "Kanban Board bauen", Status: models.StatusDoing},
			{ID: next(), Title: "Kaffee holen", Status: models.StatusDone},
		}
		_ = storage.Store.SaveTasks(tasks)
	}
}

func BoardHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	data := BoardData{}
	for _, t := range tasks {
		switch t.Status {
		case models.StatusTodo:
			data.Todo = append(data.Todo, t)
		case models.StatusDoing:
			data.Doing = append(data.Doing, t)
		case models.StatusDone:
			data.Done = append(data.Done, t)
		}
	}
	mu.Unlock()

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", 303)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Redirect(w, r, "/", 303)
		return
	}

	mu.Lock()
	tasks = append(tasks, models.Task{ID: next(), Title: title, Status: models.StatusTodo})
	_ = storage.Store.SaveTasks(tasks)
	mu.Unlock()

	http.Redirect(w, r, "/", 303)
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	status := r.FormValue("status")

	mu.Lock()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = models.Status(status)
		}
	}
	_ = storage.Store.SaveTasks(tasks)
	mu.Unlock()

	http.Redirect(w, r, "/", 303)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))

	mu.Lock()
	var newTasks []models.Task
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}
	tasks = newTasks
	_ = storage.Store.SaveTasks(tasks)
	mu.Unlock()

	http.Redirect(w, r, "/", 303)
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	tasks = []models.Task{
		{ID: next(), Title: "Go installieren", Status: models.StatusTodo},
		{ID: next(), Title: "Kanban Board bauen", Status: models.StatusDoing},
		{ID: next(), Title: "Kaffee holen", Status: models.StatusDone},
	}
	_ = storage.Store.SaveTasks(tasks)
	mu.Unlock()
	http.Redirect(w, r, "/", 303)
}
