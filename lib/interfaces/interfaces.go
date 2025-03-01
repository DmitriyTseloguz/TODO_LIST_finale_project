package interfaces

import (
	"time"
	"todo-list/lib/extensions"
)

type IReschedulable interface {
	SetTime(time.Time)
	GetTime() time.Time
	GetRepeater() extensions.RepeaterRule
}

type IRescheduler interface {
	Reschedule(IReschedulable) error
	SetBaseOnDate(time.Time)
	GetBaseOnDate() time.Time
	GetNextDate(IReschedulable) time.Time
	GetType() string
	GetOptions() string
}
