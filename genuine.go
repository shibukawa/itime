package itime

import "time"

type GenuineTime struct {}

func (GenuineTime) Now() time.Time {
	return time.Now()
}

func (GenuineTime) NewTimer(d time.Duration) Timer {
	return &GenuineTimer{timer: time.NewTimer(d)}
}

func (GenuineTime) AfterFunc(d time.Duration, f func()) Timer {
	return &GenuineTimer{timer: time.AfterFunc(d, f)}
}

func (GenuineTime) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (GenuineTime) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

func (GenuineTime) Sleep(d time.Duration) {
	time.Sleep(d)
}

var _ Time = &GenuineTime{}

type GenuineTimer struct {
	timer *time.Timer
}

func (t *GenuineTimer) Reset(d time.Duration) bool {
	return t.timer.Reset(d)
}

func (t *GenuineTimer) Stop() bool {
	return t.timer.Stop()
}

func (t GenuineTimer) Chan() <-chan time.Time {
	return t.timer.C
}

var _ Timer = &GenuineTimer{}



