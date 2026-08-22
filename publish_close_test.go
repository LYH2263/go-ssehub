package ssehub

import (
	"context"
	"errors"
	"testing"
)

// TestPublishAfterRoomCloseMustReturnErrClosed guards the regression where
// Room.Close() niled out r.Replay and Publish had its closed-check removed,
// so a publish into a closed room nil-derefed r.Replay.Append (crash) instead
// of returning a decidable ErrClosed.
//
// The CDN edge keeps publishing into a room that the offline script just
// closed. Upstream retry/circuit-breaker policy depends on errors.Is(err,
// ErrClosed) being stable and true. It must never see a panic.
func TestPublishAfterRoomCloseMustReturnErrClosed(t *testing.T) {
	r := NewRoom("live-100", 64, 8)
	ctx := context.Background()

	// Baseline: a live room publishes fine.
	if _, err := r.Publish(ctx, "chat", []byte("hello")); err != nil {
		t.Fatalf("publish on open room: %v", err)
	}

	// Offline script closes the room.
	r.Close()

	// CDN edge still publishing into the same room: must get ErrClosed,
	// deterministically, with no panic.
	ev, err := r.Publish(ctx, "chat", []byte("after-close"))
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("publish after close: want ErrClosed, got ev=%v err=%v", ev, err)
	}
}

// TestPublishIdempotentAcrossMultipleCloses ensures closing twice (and thus
// niling r.Replay twice) keeps Publish's verdict stable as ErrClosed. Close is
// idempotent on the flag; Publish must rely on the flag, not on replay state.
func TestPublishIdempotentAcrossMultipleCloses(t *testing.T) {
	r := NewRoom("live-200", 64, 8)
	r.Close()
	r.Close() // second close is a no-op but keeps Replay nil

	_, err := r.Publish(context.Background(), "chat", []byte("x"))
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}

// TestSubscribeAfterRoomClose mirrors the publish guarantee on the read path
// (subscribe.go already guards on r.closed before touching r.Replay). This
// documents that the read side was already correct and stays correct.
func TestSubscribeAfterRoomCloseMustReturnErrClosed(t *testing.T) {
	r := NewRoom("live-300", 64, 8)
	r.Close()
	_, err := r.Subscribe(context.Background(), "sub-1", 0, 8)
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("subscribe after close: want ErrClosed, got %v", err)
	}
}
