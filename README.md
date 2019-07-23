# itime: mockable time package

[![GoDoc](https://godoc.org/github.com/shibukawa/itime?status.svg)](https://godoc.org/github.com/shibukawa/itime)

This package provides go's time package compatible interface and its mock.

It provides ``Time`` interface that has almost same functions of ``time`` package.

And this package provides two implementation of the interface.

* ``GenuineTime``: It provides a real time functions.
* ``MockTime``: It is an mock time functions.

Mock functions doesn't consume real wall clock time. But it works as real ``time`` package.
For example, ``MockTime.Sleep()`` doesn't wait actually, but ``Timer`` and ``Ticker``'s channel works.

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/shibukawa/itime"
)

func timeTest(ctx context.Context, mock bool) {
	var t itime.Time
	if mock {
		t = itime.NewMock()
	} else {
		t = itime.New()
	}
	
	ticker := t.NewTicker(5 * time.Second)
	
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
	
	t.Sleep(time.Minute)
}
```

## License

Apache 2 license
