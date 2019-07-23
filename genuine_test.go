package itime

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenuineTime(t *testing.T) {
	genuineTime := New()
	defer genuineTime.Close()

	now := genuineTime.Now()
	assert.True(t, time.Now().Sub(now) < time.Millisecond)
}

func TestGenuineTime_NewTicker(t *testing.T) {
	genuineTime := New()

	ticker := genuineTime.NewTicker(2 * time.Millisecond)

	var counter int64

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ticker.Chan():
				atomic.AddInt64(&counter, 1)
			case <-ctx.Done():
				return
			}
		}
	}()
	// two ticks
	genuineTime.Sleep(5 * time.Millisecond)
	ticker.Stop()
	cancel()
	// no event called after stop
	genuineTime.Sleep(5 * time.Millisecond)

	finalCount := int(atomic.LoadInt64(&counter))
	t.Logf("count: %d", finalCount)
	assert.True(t, finalCount > 1)
}

func TestGenuineTime_NewTimer(t *testing.T) {
	genuineTime := New()

	timer := genuineTime.NewTimer(2 * time.Millisecond)

	var counter int64

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-timer.Chan():
				atomic.AddInt64(&counter, 1)
			case <-ctx.Done():
				return
			}
		}
	}()
	// one tick
	genuineTime.Sleep(5 * time.Millisecond)
	timer.Stop()
	cancel()
	// no event called after stop
	genuineTime.Sleep(5 * time.Millisecond)

	finalCount := int(atomic.LoadInt64(&counter))
	t.Logf("count: %d", finalCount)
	assert.True(t, finalCount == 1)
}

func TestGenuineTime_NewTimer2(t *testing.T) {
	genuineTime := New()

	timer := genuineTime.NewTimer(2 * time.Millisecond)

	var counter int64

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-timer.Chan():
				atomic.AddInt64(&counter, 1)
			case <-ctx.Done():
				return
			}
		}
	}()
	// no tick
	// stop before ring
	timer.Stop()
	// no event called after stop
	genuineTime.Sleep(5 * time.Millisecond)
	cancel()

	finalCount := int(atomic.LoadInt64(&counter))
	t.Logf("count: %d", finalCount)
	assert.True(t, finalCount == 0)
}

func TestGenuineTime_NewTimer_Reset(t *testing.T) {
	genuineTime := New()

	timer := genuineTime.NewTimer(3 * time.Millisecond)

	var counter int64

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-timer.Chan():
				atomic.AddInt64(&counter, 1)
			case <-ctx.Done():
				return
			}
		}
	}()
	genuineTime.Sleep(2 * time.Millisecond)
	timer.Reset(5 * time.Millisecond)
	genuineTime.Sleep(2 * time.Millisecond)

	// it should not be called
	finalCount := int(atomic.LoadInt64(&counter))
	t.Logf("count: %d", finalCount)
	assert.True(t, finalCount == 0)

	// stop before ring
	// no event called after stop
	genuineTime.Sleep(5 * time.Millisecond)
	cancel()

	finalCount = int(atomic.LoadInt64(&counter))
	t.Logf("count: %d", finalCount)
	assert.True(t, finalCount == 1)
}

func TestGenuineTime_After(t *testing.T) {
	genuineTime := New()

	timer := genuineTime.NewTimer(2 * time.Millisecond)

	var counter int64

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-timer.Chan():
				atomic.AddInt64(&counter, 1)
			case <-ctx.Done():
				return
			}
		}
	}()
	// one tick
	genuineTime.Sleep(5 * time.Millisecond)
	timer.Stop()
	cancel()
	// no event called after stop
	genuineTime.Sleep(5 * time.Millisecond)

	finalCount := int(atomic.LoadInt64(&counter))
	t.Logf("count: %d", finalCount)
	assert.True(t, finalCount == 1)
}

func TestGenuineTime_AfterFunc(t *testing.T) {
	genuineTime := New()

	var counter int64

	timer := genuineTime.AfterFunc(2*time.Millisecond, func() {
		atomic.AddInt64(&counter, 1)
	})

	// one tick
	genuineTime.Sleep(10 * time.Millisecond)
	timer.Stop()
	// no event called after stop
	genuineTime.Sleep(5 * time.Millisecond)

	finalCount := int(atomic.LoadInt64(&counter))
	t.Logf("count: %d", finalCount)
	assert.True(t, finalCount == 1)
}

func TestGenuineTime_AfterFunc2(t *testing.T) {
	genuineTime := New()

	var counter int64

	timer := genuineTime.AfterFunc(2*time.Millisecond, func() {
		atomic.AddInt64(&counter, 1)
	})

	// no tick
	// stop before ring
	timer.Stop()
	// no event called after stop
	genuineTime.Sleep(5 * time.Millisecond)

	finalCount := int(atomic.LoadInt64(&counter))
	t.Logf("count: %d", finalCount)
	assert.True(t, finalCount == 0)
}
