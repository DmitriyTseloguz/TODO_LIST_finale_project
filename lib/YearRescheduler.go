package planner

import (
	"todo-list/lib/interfaces"
)

type YearRescheduler struct {
	BaseRescheduler
}

func (rescheduler YearRescheduler) Reschedule(event interfaces.IReschedulable) error {
	var eventTime = event.GetTime()
	var baseOnDate = rescheduler.GetBaseOnDate()

	var shiftYear = 1
	var isSameYear = eventTime.Year() == baseOnDate.Year()

	if isSameYear || eventTime.After(baseOnDate) {
		event.SetTime(eventTime.AddDate(shiftYear, 0, 0))

		return nil
	}

	if eventTime.Year() < baseOnDate.Year() {
		eventTime = eventTime.AddDate(baseOnDate.Year()-eventTime.Year(), 0, 0)
	}

	event.SetTime(eventTime)

	return nil
}
