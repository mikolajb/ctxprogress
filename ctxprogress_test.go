package ctxprogress

import (
	"context"
	"testing"
)

func TestIfWrongValueIsStored(t *testing.T) {
	ctx := context.WithValue(context.Background(), key, "not a receiver")

	reporter := StartReporting(ctx)
	reporter.Report(1, 1)
}
