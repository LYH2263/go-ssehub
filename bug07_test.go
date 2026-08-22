package ssehub_test

import (
	"context"
	ssehub "github.com/LYH2263/go-ssehub"
	"testing"
)

func TestBug07_PublishCancelNoReplay(t *testing.T) {
	h := ssehub.NewHub()
	defer h.Close()
	r := h.OpenRoom("r1")
	before := r.PendingReplay()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := h.Fanout(ctx, "r1", "x", []byte("z")); err == nil {
		t.Fatal("want cancel")
	}
	if r.PendingReplay() != before {
		t.Fatalf("%d->%d", before, r.PendingReplay())
	}
}
