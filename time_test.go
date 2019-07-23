package itime_test

import (
	"context"
	"fmt"
	"time"

	"github.com/shibukawa/itime"
)

func ExampleNew() {
	t := itime.New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ticker := t.NewTicker(time.Millisecond)

	go func() {
		for {
			select {
			case now := <-ticker.Chan():
				fmt.Println(now)
			case <-ctx.Done():
				return
			}
		}
	}()

	t.Sleep(20 * time.Millisecond)
}

func ExampleNewMock() {
	t := itime.NewMock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ticker := t.NewTicker(time.Millisecond)

	go func() {
		for {
			select {
			case now := <-ticker.Chan():
				fmt.Println(now)
			case <-ctx.Done():
				return
			}
		}
	}()

	t.Sleep(20 * time.Millisecond)
}
