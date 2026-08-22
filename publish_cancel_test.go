package ssehub

import (
	"context"
	"testing"
)

// TestPublishCancelledContextDoesNotAppendReplay guards against the "串台" bug:
// an upstream request that is cancelled before publish completes must fail,
// and crucially must NOT have written the event into the replay ring —
// otherwise a reconnecting client replays a ghost frame and PendingReplay
// rises on a publish that should be a no-op.
func TestPublishCancelledContextDoesNotAppendReplay(t *testing.T) {
	r := NewRoom("room", 64, 8)

	before := r.PendingReplay()
	if before != 0 {
		t.Fatalf("fresh room should have empty replay, got %d", before)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // upstream already cancelled before publish

	ev, err := r.Publish(ctx, "evt", []byte("payload"))
	if err == nil {
		t.Fatalf("expected error from cancelled publish, got nil (ev=%+v)", ev)
	}
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	after := r.PendingReplay()
	if after != 0 {
		t.Fatalf("cancelled publish must not append to replay ring: want 0, got %d (ghost frame)", after)
	}

	// A subsequent live publish should still take the first ID slot — proving the
	// cancelled event never incremented the ring's next counter.
	liveEv, liveErr := r.Publish(context.Background(), "live", []byte("ok"))
	if liveErr != nil {
		t.Fatalf("live publish failed: %v", liveErr)
	}
	if liveEv.ID != 1 {
		t.Fatalf("first live event should get ID 1, got %d (cancelled publish leaked an ID)", liveEv.ID)
	}
	if r.PendingReplay() != 1 {
		t.Fatalf("after live publish replay should hold 1 event, got %d", r.PendingReplay())
	}
}

// TestFanoutCancelledContextDoesNotAppendReplay verifies the same invariant
// through the Hub.Fanout edge entrypoint, which previously swapped the caller's
// ctx for context.Background() and so bypassed cancellation entirely.
func TestFanoutCancelledContextDoesNotAppendReplay(t *testing.T) {
	h := NewHub()
	r := h.OpenRoom("room")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := h.Fanout(ctx, "room", "evt", []byte("payload"))
	if err == nil {
		t.Fatal("expected error from cancelled fanout, got nil")
	}

	if got := r.PendingReplay(); got != 0 {
		t.Fatalf("cancelled fanout must not append to replay ring: want 0, got %d", got)
	}
}
