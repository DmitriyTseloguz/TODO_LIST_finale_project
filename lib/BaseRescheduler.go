package planner

import (
	"errors"
	"strings"
	"time"
	"todo-list/lib/interfaces"
)

var availableRepeaters = []string{"d", "y"}

var dayRepeaterConditions = RepeaterConditions{
	{IsLengthEqualWith(2), errors.New("repeat rule \"D\" should have two parameters")},
	{IsParameterLessThan(401), errors.New("repeat rule \"D\" parameter can't be greater than 400")},
}

var yearRepeaterConditions = RepeaterConditions{
	{IsLengthEqualWith(1), errors.New("repeat rule \"Y\" should have one parameters")},
}

var otherRepeaterConditions = RepeaterConditions{
	{IsOneOf(availableRepeaters), errors.New("repeat rule is not supported")},
}

var ruleTypeRepeaterConditions = map[string]RepeaterConditions{
	"d": dayRepeaterConditions,
	"y": yearRepeaterConditions,
}

func DefineRescheduler(event interfaces.IReschedulable) (interfaces.IRescheduler, error) {
	var repeater = event.GetRepeater()
	var rescheduler interfaces.IRescheduler

	var repeaterData = strings.Split(repeater, " ")

	if len(repeater) == 0 {
		return rescheduler, errors.New("repeat rule can't be empty")
	}

	var repeaterType = repeaterData[0]

	var validator = ruleTypeRepeaterConditions[repeaterType]

	if validator == nil {
		validator = otherRepeaterConditions
	}

	var validationError = validator.GetError(repeater)

	if validationError != nil {
		return rescheduler, validationError
	}

	var base = BaseRescheduler{Type: repeaterType}
	base.SetBaseOnDate(time.Now())

	switch repeaterType {
	case "d":
		base.Options = repeaterData[1]
		rescheduler = &DayRescheduler{base}
	case "y":
		rescheduler = &YearRescheduler{base}
	}

	if rescheduler == nil {
		return nil, errors.New("not supported repeater")
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
