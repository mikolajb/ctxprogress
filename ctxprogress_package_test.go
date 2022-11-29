package ctxprogress_test

import (
	"context"
	"testing"

	"github.com/mikolajb/ctxprogress"
	"github.com/stretchr/testify/assert"
)

func TestNoop(t *testing.T) {
	ctx := context.Background()

	reporter := ctxprogress.StartReporting(ctx)
	reporter.Report(50, 100)
}

func TestReporting(t *testing.T) {
	ctx, receiver := ctxprogress.WithProgressReceiver(context.Background())

	current, total := receiver.Receive()
	assert.Equal(t, 0, current)
	assert.Equal(t, 0, total)

	reporter1 := ctxprogress.StartReporting(ctx)
	reporter1.Report(10, 100)

	current, total = receiver.Receive()
	assert.Equal(t, 10, current)
	assert.Equal(t, 100, total)

	reporter2 := ctxprogress.StartReporting(ctx)
	reporter2.Report(1, 10)

	current, total = receiver.Receive()
	assert.Equal(t, 11, current)
	assert.Equal(t, 110, total)
}
