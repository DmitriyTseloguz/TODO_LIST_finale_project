package model

import (
	"time"

	"todo-list/lib/extensions"
)

type Task struct {
	ID      string                `json:"id"`
	Title   string                `json:"title"`
	Date    extensions.ExtendTime `json:"date"`
	Comment string                `json:"comment"`
	Repeat  string                `json:"repeat"`
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
