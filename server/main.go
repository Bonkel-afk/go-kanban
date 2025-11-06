package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Status string

const (
	StatusTodo  Status = "todo"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status Status `json:"status"`
}

// Speicherdatei
const dataFile = "tasks.json"

var (
	tasks  []Task
	nextID = 1
	mu     sync.Mutex
)

type BoardData struct {
	Todo  []Task
	Doing []Task
	Done  []Task
}

// Template wird jetzt aus einer Datei geladen
// Pfad relativ zum Ordner, aus dem du `go run` startest
var tmpl = template.Must(template.ParseFiles("internals/web/templates/board.html"))

func main() {
	mu.Lock()
	if err := loadTasks(); err != nil {
		log.Println("Konnte tasks nicht laden:", err)
	}

	// Wenn noch nix gespeichert ist, Demo-Daten anlegen
	if len(tasks) == 0 {
		tasks = append(tasks,
			Task{ID: next(), Title: "Go installieren", Status: StatusTodo},
			Task{ID: next(), Title: "Kanban Board bauen", Status: StatusDoing},
			Task{ID: next(), Title: "Kaffee holen", Status: StatusDone},
		)
		if err := saveTasks(); err != nil {
			log.Println("Konnte Demo-tasks nicht speichern:", err)
		}
	}
	mu.Unlock()

	// Static Files (CSS/JS) z.B. in web/static/...
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("internals/web/static")),
		),
	)

	http.HandleFunc("/", boardHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/move", moveHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/reset", resetHandler)

	log.Println("Server läuft auf http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func next() int {
	id := nextID
	nextID++
	return id
}

// lädt tasks aus tasks.json (ruft man nur mit gehaltenem mu auf)
func loadTasks() error {
	f, err := os.Open(dataFile)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	var loaded []Task
	if err := json.NewDecoder(f).Decode(&loaded); err != nil {
		return err
	}

	tasks = loaded

	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	nextID = maxID + 1
	if nextID < 1 {
		nextID = 1
	}

	return nil
}

// speichert tasks in tasks.json (ruft man nur mit gehaltenem mu auf)
func saveTasks() error {
	tmp := dataFile + ".tmp"

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

	return os.Rename(tmp, dataFile)
}

func boardHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var todo, doing, done []Task
	for _, t := range tasks {
		switch t.Status {
		case StatusTodo:
			todo = append(todo, t)
		case StatusDoing:
			doing = append(doing, t)
		case StatusDone:
			done = append(done, t)
		}
	}
	data := BoardData{
		Todo:  todo,
		Doing: doing,
		Done:  done,
	}
	mu.Unlock()

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
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
	tasks = append(tasks, Task{
		ID:     next(),
		Title:  title,
		Status: StatusTodo,
	})
	if err := saveTasks(); err != nil {
		mu.Unlock()
		log.Println("Fehler beim Speichern:", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
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

	var newStatus Status
	switch statusStr {
	case "todo":
		newStatus = StatusTodo
	case "doing":
		newStatus = StatusDoing
	case "done":
		newStatus = StatusDone
	default:
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	mu.Lock()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = newStatus
			break
		}
	}
	if err := saveTasks(); err != nil {
		mu.Unlock()
		log.Println("Fehler beim Speichern:", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
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
	var newTasks []Task
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}
	tasks = newTasks

	if err := saveTasks(); err != nil {
		mu.Unlock()
		log.Println("Fehler beim Speichern:", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	tasks = nil
	nextID = 1
	tasks = append(tasks,
		Task{ID: next(), Title: "Go installieren", Status: StatusTodo},
		Task{ID: next(), Title: "Kanban Board bauen", Status: StatusDoing},
		Task{ID: next(), Title: "Kaffee holen", Status: StatusDone},
	)
	if err := saveTasks(); err != nil {
		mu.Unlock()
		log.Println("Fehler beim Speichern:", err)
		http.Error(w, "could not save", http.StatusInternalServerError)
		return
	}
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
