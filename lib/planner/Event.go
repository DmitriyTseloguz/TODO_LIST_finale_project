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

func (event Event) GetNextDate(baseOnDate time.Time) (time.Time, error) {
	var rescheduler, err = DefineRescheduler(&event)

	if err != nil {
		return time.Time{}, err
	}

	rescheduler.SetBaseOnDate(baseOnDate)

	rescheduler.Reschedule(&event)

	return event.GetTime(), nil
}
