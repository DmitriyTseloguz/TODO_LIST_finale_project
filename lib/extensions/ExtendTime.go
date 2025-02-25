package extensions

import (
	"time"
)

type ExtendTime time.Time

func (et *ExtendTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1]

	if str == "" {
		*et = ExtendTime(time.Now())
		return nil
	}

	t, err := time.Parse("20060102", str)

	if err != nil {
		return err
	}

	*et = ExtendTime(t)

	return nil
}
