package interfaces

import "time"

type IReschedulable interface {
	SetTime(time.Time)
	GetTime() time.Time
	GetRepeater() string
}
