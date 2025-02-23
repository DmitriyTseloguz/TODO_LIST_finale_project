package model

import (
	"time"

	"todo-list/lib/extensions"
)

type Task struct {
	id      int
	title   string
	date    string
	comment string
	repeat  string
}

func NewTask(id int, title, date, comment, repeat string) *Task {
	return &Task{id, title, date, comment, repeat}
}

func (task *Task) SetTime(t time.Time) {
	task.date = t.Format("20060102")
}

func (task Task) GetTime() time.Time {
	var taskEventTime, _ = time.Parse("20060102", task.date)

	return taskEventTime
}

func (task Task) GetRepeater() extensions.RepeaterRule {
	return extensions.RepeaterRule(task.repeat)
}
