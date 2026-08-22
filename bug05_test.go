package ssehub_test

import (
	"context"
	"errors"
	ssehub "github.com/LYH2263/go-ssehub"
	"github.com/LYH2263/go-ssehub/internal/backpressure"
	"testing"
)

func TestBug05_BackpressureWrapped(t *testing.T) {
	h := ssehub.NewHub()
	defer h.Close()
	r := h.OpenRoom("r1")
	if _, err := r.Subscribe(context.Background(), "s1", 0, 1); err != nil {
		t.Fatal(err)
	}
	_, _ = r.Publish(context.Background(), "a", []byte("1"))
	_, err := h.Fanout(context.Background(), "r1", "b", []byte("2"))
	if err == nil || !errors.Is(err, backpressure.ErrFull) {
		t.Fatalf("%v", err)
	}
}
