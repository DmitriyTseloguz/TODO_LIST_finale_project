package model

import (
	"time"

	"todo-list/lib/extensions"
)

type Task struct {
	ID      int                   `json:"id"`
	Title   string                `json:"title"`
	Date    extensions.ExtendTime `json:"date"`
	Comment string                `json:"comment"`
	Repeat  string                `json:"repeat"`
}

func NewTask(id int, title string, date time.Time, comment, repeat string) *Task {
	return &Task{
		ID:      id,
		Title:   title,
		Date:    extensions.ExtendTime(date),
		Comment: comment,
		Repeat:  repeat,
	}
}

func (task *Task) SetTime(t time.Time) {
	task.Date = extensions.ExtendTime(t)
}

func (task Task) GetTime() time.Time {
	return time.Time(task.Date)
}

func (task Task) GetRepeater() extensions.RepeaterRule {
	return extensions.RepeaterRule(task.Repeat)
}
