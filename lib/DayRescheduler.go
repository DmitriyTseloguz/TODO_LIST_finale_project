package planner

import (
	"strconv"
	"todo-list/lib/interfaces"
)

type DayRescheduler struct {
	BaseRescheduler
}

func (rescheduler DayRescheduler) Reschedule(event interfaces.IReschedulable) error {
	var eventTime = event.GetTime()
	var baseOnDate = rescheduler.GetBaseOnDate()

	var timeShift, _ = strconv.Atoi(rescheduler.Options)

	if eventTime.Before(baseOnDate) {
		var dayDifference = int(baseOnDate.Sub(eventTime).Hours() / 24)

		if dayDifference > timeShift {
			timeShift *= (dayDifference / timeShift) + 1
		}
	}

	event.SetTime(eventTime.AddDate(0, 0, timeShift))

	return nil
}
