package timemodule

import "time"

type Clock interface {
	Now() time.Time
}

type clockImpl struct{}

func (*clockImpl) Now() time.Time {
	return time.Now().UTC()
}

func NewClock() Clock {
	return &clockImpl{}
}
