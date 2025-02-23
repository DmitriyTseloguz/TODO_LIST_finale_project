package planner

import (
	"slices"
	"strconv"
	"strings"
)

type RepeatCondition struct {
	cheker          func(string) bool
	validationError error
}

type RepeaterConditions []RepeatCondition

func (conditions RepeaterConditions) IsEveryConditionOk(repeater string) bool {
	for _, condition := range conditions {
		if !condition.cheker(repeater) {
			return false
		}
	}

	return true
}

func (conditions RepeaterConditions) GetError(repeater string) error {
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

func IsRuleEqualWith(rule string) func(string) bool {
	return func(repeater string) bool {
		return strings.Split(repeater, " ")[0] == rule
	}
}

func IsLengthEqualWith(length int) func(string) bool {
	return func(repeater string) bool {
		var repeaterLength = len(strings.Split(repeater, " "))

		return repeaterLength == length
	}
}

func IsParameterLessThan(number int) func(string) bool {
	return func(repeater string) bool {
		var repeaterData = strings.Split(repeater, " ")
		var value, _ = strconv.Atoi(repeaterData[1])

		return value < number
	}
}

func IsOneOf(values []string) func(string) bool {
	return func(repeater string) bool {
		var repeaterType = strings.Split(repeater, " ")[0]

		return slices.Contains(values, repeaterType)
	}
}
