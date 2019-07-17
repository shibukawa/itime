package itime

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestNewMock(t *testing.T) {
	now := time.Now()
	mock := NewMock()

	assert.True(t, mock.Now().Sub(now) < 10 * time.Millisecond)
}

func TestNewMockWith(t *testing.T) {
	now := time.Date(2017, time.July, 17, 8, 44, 0, 0, time.UTC)
	mock := NewMockWith(now)

	assert.True(t, mock.Now().Equal(now))
}
