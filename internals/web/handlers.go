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
	Tasks []models.Task
	Tmpl  = template.Must(template.ParseFiles("internals/web/templates/board.html"))
)

func BoardHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var todo, doing, done []models.Task
	for _, t := range Tasks {
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

	if err := Tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	id := len(Tasks) + 1
	newTask := models.Task{ID: id, Title: title, Status: models.StatusTodo}
	Tasks = append(Tasks, newTask)
	storage.SaveTasks("tasks.json", Tasks)
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func MoveTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	statusStr := r.FormValue("status")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	for i := range Tasks {
		if Tasks[i].ID == id {
			switch statusStr {
			case "todo":
				Tasks[i].Status = models.StatusTodo
			case "doing":
				Tasks[i].Status = models.StatusDoing
			case "done":
				Tasks[i].Status = models.StatusDone
			}
			break
		}
	}
	storage.SaveTasks("tasks.json", Tasks)
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
