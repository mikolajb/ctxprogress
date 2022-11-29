package ctxprogress

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

func StartReporting(ctx context.Context) Reporter {
	value := ctx.Value(key)

	if value == nil {
		return noop()
	}

	receiver, ok := value.(*receiver)

	if !ok {
		// executing this is very unlikely
		return noop()
	}

	return newReporter(receiver)
}

func WithProgressReceiver(ctx context.Context) (context.Context, Receiver) {
	receiver := newReceiver()

	return context.WithValue(ctx, key, receiver), receiver
}

type report struct {
	currentValue int
	total        int
}

type contextKey int

var key contextKey

type Reporter interface {
	Report(currentValue, total int)
}

type reporter struct {
	callback func(currentValue, total int)
}

func (r *reporter) Report(currentValue, total int) {
	r.callback(currentValue, total)
}

func newReporter(receiver *receiver) *reporter {
	reporterID := uuid.New().String()

	return &reporter{
		callback: func(current, all int) {
			receiver.events.Store(reporterID, &report{
				currentValue: current,
				total:        all,
			})
		},
	}
}

type Receiver interface {
	Receive() (currentValue, total int)
}

type receiver struct {
	events *sync.Map
}

func (r *receiver) Receive() (int, int) {
	currentValue, total := 0, 0

	r.events.Range(func(_, reportersProgress any) bool {
		rep := reportersProgress.(*report)
		currentValue += rep.currentValue
		total += rep.total

		return true
	})

	return currentValue, total
}

func newReceiver() *receiver {
	return &receiver{
		events: &sync.Map{},
	}
}

func noop() Reporter {
	return &reporter{
		callback: func(current, all int) {},
	}
}
