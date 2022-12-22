# Progress in context

`ctxprogress` library makes it easy to introduce progress reporting in existing applications by relying on `context.Context`.

## How it works

The concept is based on two entities: `Receiver` and `Reporter`. Reporter uses function `Report(currentValue, total int)` to submit its progress. `Receiver` uses function `Receive() (currentValue, total int)` to fetch overall progress.

E.g., if there are 3 reporters reporting their progress in the following way:

```go
reporter1.Report(1, 10)
reporter2.Report(2,5)
reporter3.Report(3,3)
```

The receiver will read the following:

```go
progress, total := receiver.Receiver()
// progress is 6 = 1 + 2 + 3
// total is 18 = 10 + 5 + 3
```

## How to introduce it into the code

Use functions:
- `ctx, receiver := ctxprogress.WithProgressReceiver(ctx)` to initialize `Receiver`
- `reporter := ctxprogress.StartReporting(ctx)` to initialize reporter

*note: Invocation of `StartReporting` creates a new reporter each time it's invoked*

### Receiver

Inject a receiver into `context`:

```go
ctx, receiver := ctxprogress.WithProgressReceiver(ctx)
```

You can immediately check progress:

```go
progress, total := receiver.Receive()
fmt.Printf("completed: %.2f", float64(progress)/float64(total)*100)
```

### Reporter

First, you need to extract `Reporter` from `context`:

```go
reporter := ctxprogress.StartReporting(ctx)
```

You can immediately start reporting:

```go
reporter.Report(123, 1234)
```

*note: If there is no receiver in `context`, noop reporter will be returned.*

## Example

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/mikolajb/ctxprogress"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    ctx, receiver := ctxprogress.WithProgressReceiver(ctx)

    wait := &sync.WaitGroup{}
    wait.Add(1)
    go func(ctx context.Context) {
        reporter := ctxprogress.StartReporting(ctx)

        for j := 0; j < 100; j++ {
            reporter.Report(j+1, 100)
            time.Sleep(100 * time.Millisecond)
        }

        wait.Done()
    }(ctx)

    go func() {
        wait.Wait()
        cancel()
    }()

    time.Sleep(50 * time.Millisecond)
    for {
        select {
        case <-ctx.Done():
            fmt.Println("DONE")
            return
        case <-time.After(500 * time.Millisecond):
            progress, total := receiver.Receive()
            fmt.Printf("%3.2f\n", float64(progress)/float64(total)*100)
        }
    }
}
```
