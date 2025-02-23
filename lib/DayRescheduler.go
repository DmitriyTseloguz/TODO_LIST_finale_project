package planner

import (
	"strconv"
	"time"
	"todo-list/lib/interfaces"
)

type DayRescheduler struct {
	BaseRescheduler
}

func (rescheduler DayRescheduler) Reschedule(event interfaces.IReschedulable) error {
	event.SetTime(rescheduler.GetNextDate(event))

	return nil
}

func (rescheduler DayRescheduler) GetNextDate(event interfaces.IReschedulable) time.Time {
	var eventTime = event.GetTime()
	var baseOnDate = rescheduler.GetBaseOnDate()

	var timeShift, _ = strconv.Atoi(rescheduler.Options)

	if eventTime.Before(baseOnDate) {
		var dayDifference = int(baseOnDate.Sub(eventTime).Hours() / 24)

		if dayDifference > timeShift {
			timeShift *= (dayDifference / timeShift) + 1
		}
	}

	return eventTime.AddDate(0, 0, timeShift)
}
