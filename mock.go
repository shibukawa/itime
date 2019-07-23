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

func (m *MockTime) Close() {
	for _, timer := range m.timers {
		timer.internalStop()
	}
	m.timers = nil
}

func (m *MockTime) NewTimer(d time.Duration) Timer {
	r := &MockTimer{
		p:       m,
		c:       make(chan time.Time),
		d:       d,
		next:    m.current.Add(d),
		oneshot: true,
	}
	m.timers = append(m.timers, r)
	return r
}

func (m *MockTime) NewTicker(d time.Duration) Ticker {
	r := &MockTimer{
		p:       m,
		c:       make(chan time.Time),
		d:       d,
		next:    m.current.Add(d),
		oneshot: false,
	}
	m.timers = append(m.timers, r)
	return r
}

func (m *MockTime) AfterFunc(d time.Duration, f func()) Timer {
	t := m.NewTimer(d).(*MockTimer)
	t.cb = f
	return t
}

func (m *MockTime) After(d time.Duration) <-chan time.Time {
	t := m.NewTimer(d).(*MockTimer)
	return t.Chan()
}

func (m *MockTime) Tick(d time.Duration) <-chan time.Time {
	t := m.NewTicker(d)
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
	if t.Before(m.current) || len(m.timers) == 0 {
		m.current = t
		return
	}
	for {
		timers := make([]*MockTimer, 0, len(m.timers))
		for _, timer := range m.timers {
			if timer.next.Before(t) || timer.next.Equal(t) {
				timers = append(timers, timer)
			}
		}
		if len(timers) == 0 {
			break
		}
		sort.Slice(timers, func(i, j int) bool {
			return timers[i].next.Before(timers[j].next)
		})
		timer := timers[0]
		m.current = timer.next
		if processTimer {
			success := false
			for i := 0; i < 10; i++ {
				select {
				case timer.c <- timer.next:
					success = true
				default:
				}
				time.Sleep(time.Millisecond)
				if success {
					break
				}
			}
			if timer.cb != nil {
				timer.cb()
			}
		}
		if timer.oneshot {
			timer.Stop()
		} else {
			timer.next = timer.next.Add(timer.d)
		}
	}

	m.current = t
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
	p       *MockTime
	c       chan time.Time
	d       time.Duration
	next    time.Time
	cb      func()
	oneshot bool
	closed  bool
}

func (m *MockTimer) Reset(d time.Duration) bool {
	if m.closed {
		return false
	}
	m.d = d
	m.next = m.p.Now().Add(d)
	return true
}

func (m *MockTimer) internalStop() bool {
	if m.closed {
		return false
	}
	m.closed = true
	return true
}

func (m *MockTimer) Stop() bool {
	ok := m.internalStop()
	if !ok {
		return false
	}
	timers := make([]*MockTimer, 0, len(m.p.timers)-1)
	for _, timer := range m.p.timers {
		if timer != m {
			timers = append(timers, timer)
		}
	}
	m.p.timers = timers
	return true
}

func (m *MockTimer) Chan() <-chan time.Time {
	return m.c
}

var _ Timer = &MockTimer{}
