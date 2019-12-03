package itime

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestSequence_Event(t *testing.T) {
	r, w := io.Pipe()

	mt := NewMock()

	err := NewSequence(Option{
		Time: mt,
	}).
		Wait(2 * time.Second).
		Event(func() {
			io.WriteString(w, "Hello World")
			w.Close()
		}).
		Wait(2 * time.Second).
		Do(func() {
			startAt := mt.Now()
			c, err := ioutil.ReadAll(r)
			readAt := mt.Now()
			assert.NoError(t, err)
			assert.Equal(t, "Hello World", string(c))
			assert.Equal(t, readAt.Sub(startAt), 2*time.Second)
		})

	assert.NoError(t, err)
}

func TestMockSequence(t *testing.T) {
	mt := NewMock()
	defer mt.Close()

	var ctx context.Context

	err := NewSequence(Option{
		Time: mt,
	}).
		Wait(time.Second).
		Timeout(&ctx).
		Wait(2 * time.Second).
		Do(func() {
			log.Println(ctx)
			assert.NoError(t, ctx.Err())
			mt.Sleep(2 * time.Second)
			assert.Equal(t, context.DeadlineExceeded, ctx.Err())
		})

	assert.NoError(t, err)
}
