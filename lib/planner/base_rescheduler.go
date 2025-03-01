package planner

import (
	"errors"
	"time"
	"todo-list/lib/interfaces"
)

func DefineRescheduler(event interfaces.IReschedulable) (interfaces.IRescheduler, error) {
	var repeater = event.GetRepeater()
	var rescheduler interfaces.IRescheduler

	var repeaterType = repeater.GetType()

	var base = BaseRescheduler{Type: repeaterType}
	base.SetBaseOnDate(time.Now())

	switch repeaterType {
	case "d":
		base.Options = repeater.GetOptions()
		rescheduler = &DayRescheduler{base}
	case "y":
		rescheduler = &YearRescheduler{base}
	default:
		return rescheduler, errors.New("not supported repeater")
	}

	return rescheduler, nil
}

type BaseRescheduler struct {
	Type       string
	Options    string
	baseOnDate time.Time
}

func (rescheduler BaseRescheduler) GetBaseOnDate() time.Time {
	return rescheduler.baseOnDate
}

func (rescheduler *BaseRescheduler) SetBaseOnDate(date time.Time) {
	rescheduler.baseOnDate = date
}

func (rescheduler BaseRescheduler) GetType() string {
	return rescheduler.Type
}

func (rescheduler BaseRescheduler) GetOptions() string {
	return rescheduler.Options
}
