package ctxprogress_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mikolajb/ctxprogress"
)

func ExampleStartReporting() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx, receiver := ctxprogress.WithProgressReceiver(ctx)

	wait := &sync.WaitGroup{}
	wait.Add(1)
	go func(ctx context.Context) {
		r := ctxprogress.StartReporting(ctx)

		for j := 0; j < 100; j++ {
			r.Report(j+1, 100)
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
			progress, all := receiver.Receive()
			fmt.Printf("%3.2f\n", float64(progress)/float64(all)*100)
		}
	}

	// Output: 6.00
	// 11.00
	// 16.00
	// 21.00
	// 26.00
	// 31.00
	// 36.00
	// 41.00
	// 46.00
	// 51.00
	// 56.00
	// 61.00
	// 66.00
	// 71.00
	// 76.00
	// 81.00
	// 86.00
	// 91.00
	// 96.00
	// DONE
}
