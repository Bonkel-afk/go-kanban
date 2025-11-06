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

var (
	mu    sync.Mutex
	tmpl  = template.Must(template.ParseFiles("internals/web/templates/board.html"))
	tasks []models.Task
)

// InitTasks wird einmal im main() aufgerufen, nachdem der Storage gebaut wurde.
func InitTasks(store storage.Storage) {
	mu.Lock()
	defer mu.Unlock()

	loaded, err := store.LoadTasks()
	if err != nil {
		log.Println("⚠️ Konnte Tasks nicht laden:", err)
		tasks = []models.Task{}
	} else {
		tasks = loaded
	}

	log.Printf("🔎 InitTasks: %d Tasks geladen\n", len(tasks))

	// Wenn noch keine Tasks da sind, Demo-Daten anlegen
	if len(tasks) == 0 {
		log.Println("ℹ️ Keine Tasks gefunden – Demo-Daten werden angelegt")
		tasks = []models.Task{
			{ID: 1, Title: "Go installieren", Status: models.StatusTodo},
			{ID: 2, Title: "Kanban Board bauen", Status: models.StatusDoing},
			{ID: 3, Title: "Kaffee holen", Status: models.StatusDone},
		}
		if err := store.SaveTasks(tasks); err != nil {
			log.Println("⚠️ Konnte Demo-Tasks nicht speichern:", err)
		}
	}
}

func BoardHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	mu.Lock()
	defer mu.Unlock()

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
	}{
		Todo:  todo,
		Doing: doing,
		Done:  done,
	}

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
	defer mu.Unlock()

	// neue ID: max(ID)+1
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newTask := models.Task{
		ID:     maxID + 1,
		Title:  title,
		Status: models.StatusTodo,
	}
	tasks = append(tasks, newTask)

	if err := store.SaveTasks(tasks); err != nil {
		log.Println("Fehler beim Speichern nach Add:", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}

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

	idStr := r.FormValue("id")
	statusStr := r.FormValue("status")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	newStatus := models.Status(statusStr)

	mu.Lock()
	defer mu.Unlock()

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = newStatus
			break
		}
	}

	if err := store.SaveTasks(tasks); err != nil {
		log.Println("Fehler beim Speichern nach Move:", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, store storage.Storage) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var newTasks []models.Task
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}
	tasks = newTasks

	if err := store.SaveTasks(tasks); err != nil {
		log.Println("Fehler beim Speichern nach Delete:", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
