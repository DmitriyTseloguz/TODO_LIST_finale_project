package extensions

import (
	"errors"
	"strconv"
)

var availableRepeaters = RepeaterTypes{
	{"d", true},
	{"y", false},
}

var dayRepeaterConditions = RepeaterConditions{
	{IsLengthEqualWith(2), errors.New("repeat rule \"D\" should have two parameters")},
	{IsParameterLessThan(401), errors.New("repeat rule \"D\" parameter can't be greater than 400")},
}

var yearRepeaterConditions = RepeaterConditions{
	{IsLengthEqualWith(1), errors.New("repeat rule \"Y\" should have one parameters")},
}

var commonRepeaterConditions = RepeaterConditions{
	{IsOneOf(availableRepeaters), errors.New("repeat rule is not supported")},
}

var ruleTypeRepeaterConditions = map[string]RepeaterConditions{
	"d": dayRepeaterConditions,
	"y": yearRepeaterConditions,
}

type RepeatCondition struct {
	cheker          func(RepeaterRule) bool
	validationError error
}

type RepeaterConditions []RepeatCondition

func (conditions RepeaterConditions) IsEveryConditionOk(repeater RepeaterRule) bool {
	for _, condition := range conditions {
		if !condition.cheker(repeater) {
			return false
		}
	}

	return true
}

func (conditions RepeaterConditions) GetError(repeater RepeaterRule) error {
	for _, condition := range conditions {
		if !condition.cheker(repeater) {
			return condition.validationError
		}

	}

	return nil
}

func IsNotEmpty(repeater string) bool {
	return repeater != ""
}

func IsRuleEqualWith(rule string) func(RepeaterRule) bool {
	return func(repeater RepeaterRule) bool {
		return repeater.Separate()[0] == rule
	}
}

func IsLengthEqualWith(length int) func(RepeaterRule) bool {
	return func(repeater RepeaterRule) bool {
		var repeaterLength = len(repeater.Separate())

		return repeaterLength == length
	}
}

func IsParameterLessThan(number int) func(RepeaterRule) bool {
	return func(repeater RepeaterRule) bool {
		var repeaterData = repeater.Separate()
		var value, _ = strconv.Atoi(repeaterData[1])

		return value < number
	}
}

func IsOneOf(values RepeaterTypes) func(RepeaterRule) bool {
	return func(repeater RepeaterRule) bool {
		var repeaterType = repeater.GetType()

		return values.Some(func(value RepeaterType) bool {
			return value.display == repeaterType
		})
	}
}
