package jenkins

import (
	"context"
	"testing"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func TestEpochMsToTime(t *testing.T) {
	ctx := context.Background()

	got, err := epochMsToTime(ctx, &transform.TransformData{Value: int64(1_700_000_000_000)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts, ok := got.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", got)
	}
	if !ts.Equal(time.UnixMilli(1_700_000_000_000).UTC()) {
		t.Fatalf("unexpected time: %v", ts)
	}

	got, err = epochMsToTime(ctx, &transform.TransformData{Value: int64(0)})
	if err != nil {
		t.Fatalf("unexpected error for zero: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for non-positive timestamp, got %v", got)
	}

	got, err = epochMsToTime(ctx, &transform.TransformData{Value: nil})
	if err != nil {
		t.Fatalf("unexpected error for nil: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for nil value, got %v", got)
	}
}
