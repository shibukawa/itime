package itime

import "time"

// Ticker is an interface of oneshot timer.
type Timer interface {
	Reset(d time.Duration) bool
	Stop() bool
	Chan() <-chan time.Time
}

// Ticker is an interface of interval timer.
type Ticker interface {
	Stop() bool
	Chan() <-chan time.Time
}

// Time is an interface of entrypoint of this package.
//
// It provides compatible features of go's time package as much as possible.
type Time interface {
	Now() time.Time
	NewTimer(d time.Duration) Timer
	NewTicker(d time.Duration) Ticker
	AfterFunc(d time.Duration, f func()) Timer
	After(d time.Duration) <-chan time.Time
	Tick(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Close()
}
