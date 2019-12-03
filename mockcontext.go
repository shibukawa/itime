package itime

import (
	"context"
	"runtime"
	"sync"
	"time"
)

type mockContext struct {
	time     *MockTime
	deadline time.Time
	ch       chan struct{}
	err      error
	parent   context.Context
}

func newMockContext(t *MockTime, parent context.Context, d time.Time) (*mockContext, context.CancelFunc) {
	timer := t.NewTimer(d.Sub(t.Now()))
	var once sync.Once

	c := &mockContext{
		time:     t,
		deadline: d,
		ch:       make(chan struct{}),
		parent:   parent,
	}
	closeChan := func(err error) {
		once.Do(func() {
			c.err = err
		})
	}
	cancel := func() {
		closeChan(context.Canceled)
	}
	go func() {
		select {
		case <-parent.Done():
			closeChan(parent.Err())
		case <-timer.Chan():
			closeChan(context.DeadlineExceeded)
		}
		timer.Stop()
	}()
	runtime.Gosched()
	return c, cancel
}

func (c mockContext) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c mockContext) Done() <-chan struct{} {
	return c.ch
}

func (c mockContext) Err() error {
	if c.err == nil {
		return c.parent.Err()
	}
	return c.err
}

func (c mockContext) Value(key interface{}) interface{} {
	return c.parent.Value(key)
}

var _ context.Context = &mockContext{}
