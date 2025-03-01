package planner

import (
	"time"
	"todo-list/lib/interfaces"
)

type YearRescheduler struct {
	BaseRescheduler
}

func (rescheduler YearRescheduler) Reschedule(event interfaces.IReschedulable) error {
	event.SetTime(rescheduler.GetNextDate(event))

	return nil
}

func (rescheduler YearRescheduler) GetNextDate(event interfaces.IReschedulable) time.Time {
	var eventTime = event.GetTime()
	var baseOnDate = rescheduler.GetBaseOnDate()

	var shiftYear = 1
	var isSameYear = eventTime.Year() == baseOnDate.Year()

	if isSameYear || eventTime.After(baseOnDate) {
		return eventTime.AddDate(shiftYear, 0, 0)
	}

	if eventTime.Year() < baseOnDate.Year() {
		eventTime = eventTime.AddDate(baseOnDate.Year()-eventTime.Year(), 0, 0)
	}

	return eventTime
}
