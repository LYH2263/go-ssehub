package ssehub

import (
	"context"
	"errors"
	"testing"
)

// TestSubscribeMissingEncoder reproduces the pressure-test scenario where the
// Encoder is nilled out to forcibly stop encoding, yet a client still calls
// Subscribe. Before the fix the symptoms were:
//  1. first Subscribe "直接炸" — nil-encoder deref during replay panicked;
//  2. re-subscribe with the same subID returned ErrConflict (half-slot left);
//  3. HasSubscriber still reported true;
//  4. other subscriptions/flows could not progress.
//
// After the fix: a missing encoder yields a decidable ErrInvalid error and
// leaves no half-registered slot behind.
func TestSubscribeMissingEncoder(t *testing.T) {
	t.Run("fails with decidable ErrInvalid and no panic", func(t *testing.T) {
		room := NewRoom("orders", 64, 8)
		// Seed replay so the old code would reach r.Enc.Encode during Subscribe.
		ctx := context.Background()
		if _, err := room.Publish(ctx, "created", []byte("ord-1")); err != nil {
			t.Fatalf("seed publish: %v", err)
		}
		// Stress lever: nil out the encoder to forcibly stop encoding.
		room.SetEncoder(nil)

		ch, err := room.Subscribe(ctx, "sub-1", 0, 8)
		if err == nil {
			t.Fatalf("expected error when encoder is nil, got channel %v", ch)
		}
		if !errors.Is(err, ErrInvalid) {
			t.Fatalf("error must wrap ErrInvalid for decidable handling, got %v", err)
		}
	})

	t.Run("no residual slot — same subID re-subscribes cleanly once encoder restored", func(t *testing.T) {
		room := NewRoom("orders", 64, 8)
		ctx := context.Background()
		room.SetEncoder(nil)

		if _, err := room.Subscribe(ctx, "sub-1", 0, 8); !errors.Is(err, ErrInvalid) {
			t.Fatalf("first subscribe with nil encoder: want ErrInvalid, got %v", err)
		}

		// Symptom 3: HasSubscriber must be false — the failed subscribe must not
		// leave a half-slot in the roster.
		if room.HasSubscriber("sub-1") {
			t.Fatalf("HasSubscriber=true after failed subscribe; roster has residual slot")
		}
		// Symptom 2: re-subscribing the same subID must not hit ErrConflict.
		if _, err := room.Subscribe(ctx, "sub-1", 0, 8); !errors.Is(err, ErrInvalid) {
			t.Fatalf("re-subscribe with same subID: want ErrInvalid (not Conflict), got %v", err)
		}

		// Restore the encoder; the same subID must now succeed and be observable.
		room.SetEncoder(NewDefaultEncoder())
		ch, err := room.Subscribe(ctx, "sub-1", 0, 8)
		if err != nil {
			t.Fatalf("subscribe after restoring encoder: %v", err)
		}
		if !room.HasSubscriber("sub-1") {
			t.Fatalf("HasSubscriber=false after successful subscribe")
		}
		// Symptom 4: other flows must still subscribe.
		if _, err := room.Subscribe(ctx, "sub-2", 0, 8); err != nil {
			t.Fatalf("other subscription blocked: %v", err)
		}
		if got := len(room.ListSubscribers()); got != 2 {
			t.Fatalf("expected 2 subscribers, got %d", got)
		}
		// Drain the channels so they don't block goroutine teardown.
		room.Unsubscribe("sub-1")
		room.Unsubscribe("sub-2")
		_ = ch
	})

	t.Run("empty replay does not mask the missing encoder", func(t *testing.T) {
		// With no replay events the old code never deref'd Enc, so it registered
		// the slot and returned a live channel — silently subscribing against a
		// room that can't encode. The encoder check must still fire.
		room := NewRoom("orders", 64, 8)
		room.SetEncoder(nil)
		ctx := context.Background()

		if _, err := room.Subscribe(ctx, "sub-x", 0, 8); !errors.Is(err, ErrInvalid) {
			t.Fatalf("empty-replay subscribe with nil encoder: want ErrInvalid, got %v", err)
		}
		if room.HasSubscriber("sub-x") {
			t.Fatalf("residual slot left behind on empty replay")
		}
	})
}
