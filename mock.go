package itime

import "time"

type timer struct {
	c       chan time.Time
	next    time.Time
	cb      func()
	oneshot bool
}

type MockTime struct {
	timers  map[*MockTimer]*timer
	current time.Time
}

func (m MockTime) Now() time.Time {
	return m.current
}

func isClosed(ch chan time.Time) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

func (m *MockTime) Close() error {
	return nil
}

func (MockTime) NewTimer(d time.Duration) Timer {
	panic("implement me")
}

func (MockTime) AfterFunc(d time.Duration, f func()) Timer {
	panic("implement me")
}

func (m *MockTime) After(d time.Duration) <-chan time.Time {
	t := m.NewTimer(d)
	m.timers[t.(*MockTimer)].oneshot = true
	return t.Chan()
}

func (m *MockTime) Tick(d time.Duration) <-chan time.Time {
	t := m.NewTimer(d)
	return t.Chan()
}

func (m *MockTime) Sleep(d time.Duration) {
	m.Advance(d, true)
}

func (m *MockTime) Advance(d time.Duration, processTimer bool) {
	newCurrent := m.current.Add(d)
	m.Set(newCurrent, processTimer)
}

func (m *MockTime) Set(t time.Time, processTimer bool) {
	panic("implement me")
}

var _ Time = &MockTime{}

func NewMock() *MockTime {
	return &MockTime{
		timers:  make(map[*MockTimer]*timer),
		current: time.Now(),
	}
}

func NewMockWith(t time.Time) *MockTime {
	return &MockTime{
		timers:  make(map[*MockTimer]*timer),
		current: t,
	}
}

type MockTimer struct {
	p *MockTime
	d time.Duration
	c chan time.Time
}

func (m *MockTimer) Reset(d time.Duration) bool {
	if isClosed(m.c) {
		return false
	}
	m.d = d
	m.p.timers[m].next = m.p.Now().Add(d)
	return true
}

func (m *MockTimer) Stop() bool {
	if isClosed(m.c) {
		return false
	}
	close(m.c)
	return true
}

func (m MockTimer) Chan() <-chan time.Time {
	return m.c
}

func (m *MockTimer) AdvanceToNext() {
}

var _ Timer = &MockTimer{}
