package itime

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMock(t *testing.T) {
	now := time.Now()
	mock := NewMock()
	defer mock.Close()

	assert.True(t, mock.Now().Sub(now) < 10*time.Millisecond)
}

func TestMockTime_Close(t *testing.T) {
	mock := NewMock()

	timer := mock.NewTimer(time.Minute)

	// Close closes all timers
	mock.Close()

	// Channel is already closed
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	success := false
	select {
	case <-timer.Chan():
		success = true
	case <-ctx.Done():
	}

	assert.False(t, success)

	// Stop() returns false when it is already closed
	assert.False(t, timer.Stop())

}

func TestNewMockWith(t *testing.T) {
	now := time.Date(2017, time.July, 17, 8, 44, 0, 0, time.UTC)
	mock := NewMockWith(now)
	defer mock.Close()

	assert.True(t, mock.Now().Equal(now))
}

func TestMockTime_NewTicker(t *testing.T) {
	now := time.Date(2017, time.July, 17, 8, 44, 0, 0, time.UTC)
	mock := NewMockWith(now)
	defer mock.Close()

	ticker := mock.NewTicker(time.Second)
	defer ticker.Stop()
	var counter int64 = 0
	go func() {
		for {
			now, ok := <-ticker.Chan()
			if !ok {
				break
			}
			atomic.AddInt64(&counter, 1)
			t.Log(now)
		}
	}()
	mock.Advance(10*time.Second, true)

	finalCount := int(atomic.LoadInt64(&counter))
	assert.Equal(t, 10, finalCount)

	assert.True(t, ticker.Stop())
	assert.False(t, ticker.Stop())
	mock.Advance(10*time.Second, true)

	finalCount = int(atomic.LoadInt64(&counter))
	assert.Equal(t, 10, finalCount)
}

func TestMockTime_AfterFunc(t *testing.T) {
	now := time.Date(2017, time.July, 17, 8, 44, 0, 0, time.UTC)
	mock := NewMockWith(now)
	defer mock.Close()

	var counter int64 = 0

	// When call Stop(), callback function wouldn't be fired
	timer1 := mock.AfterFunc(2*time.Second, func() {
		atomic.AddInt64(&counter, 1)
	})
	mock.Advance(time.Second, true)
	timer1.Stop()
	mock.Advance(2*time.Second, true)

	finalCount := int(atomic.LoadInt64(&counter))
	assert.Equal(t, 0, finalCount)

	timer2 := mock.AfterFunc(time.Second, func() {
		atomic.AddInt64(&counter, 1)
	})
	defer timer2.Stop()
	mock.Advance(2*time.Second, true)

	finalCount = int(atomic.LoadInt64(&counter))
	assert.Equal(t, 1, finalCount)
}

func TestMockTime_MultipleTimer(t *testing.T) {
	mock := NewMock()

	events := make(chan string, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timer1 := mock.After(time.Second * 7)
	go func() {
		for {
			select {
			case <-timer1:
				events <- "timer1(oneshot)"
			case <-ctx.Done():
				return
			}
		}
	}()

	timer2 := mock.Tick(time.Second * 3)
	go func() {
		for {
			select {
			case <-timer2:
				events <- "timer2(tick)"
			case <-ctx.Done():
				return
			}
		}
	}()

	mock.Sleep(10 * time.Second)

	assert.Equal(t, "timer2(tick)", <-events)
	assert.Equal(t, "timer2(tick)", <-events)
	assert.Equal(t, "timer1(oneshot)", <-events)
	assert.Equal(t, "timer2(tick)", <-events)

	mock.Close()

}

func TestMockTimer_Reset(t *testing.T) {
	now := time.Date(2017, time.July, 17, 8, 44, 0, 0, time.UTC)
	mock := NewMockWith(now)
	defer mock.Close()

	timer := mock.NewTimer(time.Second).(*MockTimer)
	assert.True(t, timer.Reset(time.Minute))

	wait := make(chan string)
	go func() {
		<-timer.Chan()
		wait <- "ring"
	}()
	mock.Sleep(2 * time.Minute)

	message, ok := <-wait
	assert.Equal(t, "ring", message)
	assert.True(t, ok)

	// Reset return false after stop
	assert.False(t, timer.Reset(3*time.Minute))
}
