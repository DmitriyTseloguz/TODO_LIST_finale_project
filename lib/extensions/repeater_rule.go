package extensions

import (
	"errors"
	"strings"
)

type RepeaterRule string

func (repeater RepeaterRule) Validate() error {
	var repeaterData = repeater.Separate()

	if len(repeater) == 0 {
		return errors.New("repeat rule can't be empty")
	}

	var repeaterType = repeaterData[0]

	var validator = ruleTypeRepeaterConditions[repeaterType]

	if validator == nil {
		validator = commonRepeaterConditions
	}

	return validator.GetError(repeater)
}

func (repeater RepeaterRule) Separate() []string {
	return strings.Split(string(repeater), " ")
}

func (repeater RepeaterRule) GetType() string {
	return repeater.Separate()[0]
}

func (repeater RepeaterRule) GetOptions() string {
	return repeater.Separate()[1]
}

func (repeater RepeaterRule) HasOptions() bool {
	return len(repeater.Separate()) > 1
}
