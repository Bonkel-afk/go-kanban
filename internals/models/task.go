package models

type Status string

const (
	StatusTodo  Status = "todo"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

type Task struct {
	ID     int    `json:"id" bson:"id"`
	Title  string `json:"title" bson:"title"`
	Status Status `json:"status" bson:"status"`
}
