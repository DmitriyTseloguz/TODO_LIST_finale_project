package planner

import (
	"time"
	"todo-list/lib/extensions"
)

type Event struct {
	date     time.Time
	repeater extensions.RepeaterRule
}

func NewEvent(date time.Time, repeater extensions.RepeaterRule) *Event {
	return &Event{date, repeater}
}

func (event *Event) SetTime(t time.Time) {
	event.date = t
}

func (event Event) GetTime() time.Time {
	return event.date
}

func (event Event) GetRepeater() extensions.RepeaterRule {
	return event.repeater
}
