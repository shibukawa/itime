package itime

import (
	"context"
	"errors"
	"time"
)

type Sequence struct {
	time            *MockTime
	scenarioTimeout time.Duration
	verbose         bool
	sequence        []func()
	current         time.Time
}

type Option struct {
	Time            *MockTime
	ScenarioTimeout time.Duration
	Verbose         bool
}

func NewSequence(opt Option) *Sequence {
	if opt.Time == nil {
		opt.Time = NewMock()
	}
	if opt.ScenarioTimeout == 0 {
		opt.ScenarioTimeout = time.Second * 3
	}
	return &Sequence{
		time:            opt.Time,
		verbose:         opt.Verbose,
		scenarioTimeout: opt.ScenarioTimeout,
		current:         opt.Time.Now(),
	}
}

func (s *Sequence) Wait(d time.Duration) *Sequence {
	s.sequence = append(s.sequence, func() {
		s.time.Advance(d, true)
	})
	s.current = s.current.Add(d)
	return s
}

func (s *Sequence) Timeout(ctx *context.Context) *Sequence {
	*ctx, _ = s.time.WithDeadline(context.Background(), s.current)
	return s
}

func (s *Sequence) Event(callback func()) *Sequence {
	s.sequence = append(s.sequence, func() {
		callback()
	})
	return s
}

func (s *Sequence) Do(testfunc func()) error {
	finish := make(chan error)
	go func() {
		testfunc()
		finish <- nil
	}()
	for _, step := range s.sequence {
		step()
	}
	timer := time.NewTimer(s.scenarioTimeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return errors.New("test scenario is timed out")
	case err := <-finish:
		return err
	}
}
