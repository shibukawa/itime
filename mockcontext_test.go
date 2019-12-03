package itime

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMockContext_Cancel(t *testing.T) {
	mt := NewMock()
	ctx, cancel := mt.WithTimeout(context.Background(), time.Second)
	deadline, ok := ctx.Deadline()
	assert.Equal(t, mt.Now().Add(time.Second), deadline)
	assert.True(t, ok)
	assert.Nil(t, ctx.Err())
	cancel()
	assert.Equal(t, ctx.Err(), context.Canceled)
}

func TestMockContext_Timeout(t *testing.T) {
	mt := NewMock()
	ctx, cancel := mt.WithTimeout(context.Background(), time.Second)
	deadline, ok := ctx.Deadline()
	assert.Equal(t, mt.Now().Add(time.Second), deadline)
	assert.True(t, ok)
	assert.Nil(t, ctx.Err())

	mt.Advance(time.Minute, true)
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	cancel()
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}

func TestMockContext_Deadline(t *testing.T) {
	mt := NewMock()
	d := mt.Now().Add(time.Second)
	ctx, cancel := mt.WithDeadline(context.Background(), d)
	deadline, ok := ctx.Deadline()
	assert.Equal(t, mt.Now().Add(time.Second), deadline)
	assert.True(t, ok)
	assert.Nil(t, ctx.Err())

	mt.Advance(time.Minute, true)
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	cancel()
	// If timeout first, Err() keeps first error even if cancel called.
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}

func TestMockContext_CancelByParent(t *testing.T) {
	mt := NewMock()
	parent, pcancel := context.WithCancel(context.Background())
	ctx, cancel := mt.WithTimeout(parent, time.Second)
	deadline, ok := ctx.Deadline()
	assert.Equal(t, mt.Now().Add(time.Second), deadline)
	assert.True(t, ok)
	assert.Nil(t, ctx.Err())

	pcancel()
	assert.Equal(t, context.Canceled, ctx.Err())
	mt.Advance(time.Minute, true)
	// If parent's context fulfilled first, Err keeps parent's Err
	// even if timeout happens or cancel function is called after that.
	assert.Equal(t, context.Canceled, ctx.Err())
	cancel()
	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestMockContext_TimeoutByParent(t *testing.T) {
	mt := NewMock()
	parent, _ := mt.WithTimeout(context.Background(), time.Second)
	ctx, cancel := mt.WithTimeout(parent, 10*time.Second)
	deadline, ok := ctx.Deadline()
	assert.Equal(t, mt.Now().Add(10*time.Second), deadline)
	assert.True(t, ok)
	assert.Nil(t, ctx.Err())

	mt.Advance(time.Second*2, true)
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	// If parent's context fulfilled first, Err keeps parent's Err
	// even if timeout happens or cancel function is called after that.
	mt.Advance(time.Minute, true)
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	cancel()
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}

func TestMockContext_Value(t *testing.T) {
	mt := NewMock()
	var key = "key"
	parent := context.WithValue(context.Background(), key, "value")
	ctx, cancel := mt.WithTimeout(parent, 10*time.Second)
	defer cancel()

	assert.Equal(t, "value", ctx.Value(key))
}
