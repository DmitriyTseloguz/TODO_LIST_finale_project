package extensions

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"todo-list/lib/constants"
)

type ExtendTime time.Time

func (et *ExtendTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1]

	if str == "" {
		*et = ExtendTime(time.Now())
		return nil
	}

	t, err := time.Parse(constants.DBDateFormat, str)

	if err != nil {
		return err
	}

	*et = ExtendTime(t)

	return nil
}

func (et ExtendTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(et).Format(constants.DBDateFormat))
}

func (et *ExtendTime) Scan(value any) error {
	if value == nil {
		*et = ExtendTime(time.Now())
		return nil
	}

	dateStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	t, err := time.Parse(constants.DBDateFormat, dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}

	*et = ExtendTime(t)

	return nil
}

func (et ExtendTime) Value() (driver.Value, error) {
	return time.Time(et).Format(constants.DBDateFormat), nil
}

func (et ExtendTime) String() string {
	return time.Time(et).Format(constants.DisplayDateFormat)
}
