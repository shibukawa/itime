package itime

import (
	"sort"
	"time"
)

type MockTime struct {
	timers  []*MockTimer
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

func (m *MockTime) NewTimer(d time.Duration) Timer {
	r := &MockTimer{
		p: m,
	}
	return r
}

func (m *MockTime) AfterFunc(d time.Duration, f func()) Timer {
	t := m.NewTimer(d).(*MockTimer)
	t.oneshot = true
	t.cb = f
	return t
}

func (m *MockTime) After(d time.Duration) <-chan time.Time {
	t := m.NewTimer(d).(*MockTimer)
	t.oneshot = true
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
	if t.Before(m.current) {
		m.current = t
		return
	}
	for {
		timers := make([]*MockTimer, 0, len(m.timers))
		for _, timer := range m.timers {
			if timer.next.Before(t) {
				timers = append(timers, timer)
			}
			if len(timers) == 0 {
				break
			}
			sort.Slice(timers, func(i, j int) bool {
				return timers[i].next.Before(timers[j].next)
			})
			m.current = 
			timers[0]
			timers = timers[:0]
		}
	}


	m.current = t

	panic("implement me")
}

var _ Time = &MockTime{}

func NewMock() *MockTime {
	return &MockTime{
		current: time.Now(),
	}
}

func NewMockWith(t time.Time) *MockTime {
	return &MockTime{
		current: t,
	}
}

type MockTimer struct {
	p *MockTime
	c       chan time.Time
	d       time.Duration
	next    time.Time
	cb      func()
	oneshot bool
}

func (m *MockTimer) Reset(d time.Duration) bool {
	if isClosed(m.c) {
		return false
	}
	m.d = d
	m.next = m.p.Now().Add(d)
	return true
}

func (m *MockTimer) Stop() bool {
	if isClosed(m.c) {
		return false
	}
	close(m.c)
	return true
}

func (m *MockTimer) Chan() <-chan time.Time {
	return m.c
}

func (m *MockTimer) AdvanceToNext() {
	m.p.Set(m.next, true)
}

var _ Timer = &MockTimer{}
