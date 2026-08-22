package ssehub_test

import (
	"context"
	ssehub "github.com/LYH2263/go-ssehub"
	"testing"
)

func TestBug10_CloseFlushesBeforeClear(t *testing.T) {
	r := ssehub.NewRoom("a", 8, 8)
	if _, err := r.Publish(context.Background(), "x", []byte("z")); err != nil {
		t.Fatal(err)
	}
	if r.CloseFlushCount() == 0 {
		t.Fatal("want flush")
	}
}
