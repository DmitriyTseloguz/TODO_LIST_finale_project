package interfaces

import "time"

type IReschedulable interface {
	SetTime(time.Time)
	GetTime() time.Time
	GetRepeater() string
}

type IRescheduler interface {
	Reschedule(IReschedulable) error
	SetBaseOnDate(time.Time)
	GetBaseOnDate() time.Time
	GetType() string
	GetOptions() string
}
