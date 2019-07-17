package itime

import "time"

type Timer interface {
	Reset(d time.Duration) bool
	Stop() bool
	Chan() <-chan time.Time
}

type Time interface {
	Now() time.Time
	NewTimer(d time.Duration) Timer
	AfterFunc(d time.Duration, f func()) Timer
	After(d time.Duration) <-chan time.Time
	Tick(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
}