package ssehub

import (
	"context"
	"errors"
	"testing"

	"github.com/LYH2263/go-ssehub/internal/backpressure"
)

// TestPublishFanoutBackpressureSentinel asserts that a full subscriber send
// queue surfaces as a sentinel-matched error on every entry point the edge
// gateway uses to trigger rate limiting / circuit breaking:
//
//	errors.Is(err, backpressure.ErrFull)
//
// The fanout path used to wrap with %v, which stringified the error and
// broke errors.Is — forcing on-call to grep the message for "full".
// Both the room publish path and the hub fanout path must expose the same
// backpressure sentinel to callers.
func TestPublishFanoutBackpressureSentinel(t *testing.T) {
	ctx := context.Background()
	const room = "bp-room"

	h := NewHub()
	r := h.OpenRoom(room)
	// Tiny queue so a single un-drained subscriber saturates on the 2nd push.
	if _, err := r.Subscribe(ctx, "slow", 0, 1); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	// First publish fills the 1-cap queue (ring replay push is dropped by
	// design, the only subscriber push that can fail here is the live one).
	if _, err := h.Publish(ctx, room, "n", []byte("a")); err != nil {
		t.Fatalf("first publish should not fail, got: %v", err)
	}

	// Second publish: the slow subscriber's queue is full → ErrFull.
	_, pubErr := h.Publish(ctx, room, "n", []byte("b"))
	if pubErr == nil {
		t.Fatalf("expected backpressure error from Hub.Publish, got nil")
	}
	if !errors.Is(pubErr, backpressure.ErrFull) {
		t.Errorf("Hub.Publish: errors.Is(err, backpressure.ErrFull) = false; err = %v", pubErr)
	}

	// Fresh room to exercise the Fanout path independently.
	const froom = "bp-fanout"
	r2 := h.OpenRoom(froom)
	if _, err := r2.Subscribe(ctx, "slow", 0, 1); err != nil {
		t.Fatalf("subscribe fanout: %v", err)
	}
	if _, err := h.Fanout(ctx, froom, "n", []byte("a")); err != nil {
		t.Fatalf("first fanout should not fail, got: %v", err)
	}
	_, fanErr := h.Fanout(ctx, froom, "n", []byte("b"))
	if fanErr == nil {
		t.Fatalf("expected backpressure error from Hub.Fanout, got nil")
	}
	// The regression: Hub.Fanout must preserve the sentinel across wrapping
	// so the gateway's errors.Is check trips circuit breaking.
	if !errors.Is(fanErr, backpressure.ErrFull) {
		t.Errorf("Hub.Fanout: errors.Is(err, backpressure.ErrFull) = false; err = %v", fanErr)
	}
}
