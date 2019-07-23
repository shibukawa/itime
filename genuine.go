package itime

import "time"

// GenuineTime is an implementation of real time for Time interface
type GenuineTime struct {
	timers  []*GenuineTimer
	tickers []*GenuineTicker
}

// Close methods stops all internal timer
func (t *GenuineTime) Close() {
	for _, timer := range t.timers {
		timer.internalStop()
	}
	for _, ticker := range t.tickers {
		ticker.internalStop()
	}
	t.timers = nil
	t.tickers = nil
}

// Now returns wall clock time
func (GenuineTime) Now() time.Time {
	return time.Now()
}

// NewTimer creates new oneshot timer
func (t *GenuineTime) NewTimer(d time.Duration) Timer {
	r := &GenuineTimer{
		p: t,
		t: time.NewTimer(d),
		w: make(chan struct{}),
		c: make(chan time.Time),
	}
	t.timers = append(t.timers, r)
	go func() {
		select {
		case now := <-r.t.C:
			r.Stop()
			r.c <- now
		case <-r.w:
		}
	}()
	return r
}

// NewTicker creates interval timer
func (t *GenuineTime) NewTicker(d time.Duration) Ticker {
	r := &GenuineTicker{
		p: t,
		t: time.NewTicker(d),
	}
	t.tickers = append(t.tickers, r)
	return r
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
//
// Resulting timer is for stopping timer.
func (t *GenuineTime) AfterFunc(d time.Duration, f func()) Timer {
	r := &GenuineTimer{
		p: t,
		t: time.NewTimer(d),
		w: make(chan struct{}),
		c: make(chan time.Time),
	}
	t.timers = append(t.timers, r)
	go func() {
		select {
		case now := <-r.t.C:
			r.Stop()
			f()
			r.c <- now
		case <-r.w:
		}
	}()
	return r
}

// After is a shorthand of creating Timer instance
func (t *GenuineTime) After(d time.Duration) <-chan time.Time {
	r := t.NewTimer(d)
	return r.Chan()
}

// After is a shorthand of creating Ticker instance
func (t *GenuineTime) Tick(d time.Duration) <-chan time.Time {
	r := t.NewTicker(d)
	return r.Chan()
}

// Sleep waits for the duration to elapse
func (GenuineTime) Sleep(d time.Duration) {
	time.Sleep(d)
}

var _ Time = &GenuineTime{}

// New returns GenuineTime instance
func New() Time {
	return &GenuineTime{}
}

// GenuineTimer is an actual oneshot timer implementation of Timer interface
type GenuineTimer struct {
	p *GenuineTime
	t *time.Timer
	w chan struct{}
	c chan time.Time
}

// Reset changes timer duration
func (t *GenuineTimer) Reset(d time.Duration) bool {
	return t.t.Reset(d)
}

func (t *GenuineTimer) internalStop() bool {
	r := t.t.Stop()
	w := t.w
	t.w = nil
	if w != nil {
		close(w)
	}
	return r
}

// Stop stops timer
func (t *GenuineTimer) Stop() bool {
	ok := t.internalStop()
	if ok == true {
		timers := make([]*GenuineTimer, 0, len(t.p.timers)-1)
		for _, timer := range t.p.timers {
			if timer != t {
				timers = append(timers, timer)
			}
		}
		t.p.timers = timers
	}
	return ok
}

// Chan returns channel that sends current time
func (t GenuineTimer) Chan() <-chan time.Time {
	return t.c
}

var _ Timer = &GenuineTimer{}

// GenuineTicker is an actual interval timer implementation of Timer interface
type GenuineTicker struct {
	p *GenuineTime
	t *time.Ticker
}

func (t *GenuineTicker) internalStop() {
	t.t.Stop()
}

// Stop stops timer. But it doesn't close channel.
func (t *GenuineTicker) Stop() bool {
	t.internalStop()
	tickers := make([]*GenuineTicker, 0, len(t.p.tickers)-1)
	for _, ticker := range t.p.tickers {
		if ticker != t {
			tickers = append(tickers, ticker)
		}
	}
	t.p.tickers = tickers
	return true
}

// Chan returns channel that sends current time
func (t GenuineTicker) Chan() <-chan time.Time {
	return t.t.C
}

var _ Ticker = &GenuineTicker{}
