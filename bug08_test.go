package ssehub_test

import (
	"context"
	ssehub "github.com/LYH2263/go-ssehub"
	"testing"
	"time"
)

func TestBug08_WaitFrameHonorsCancel(t *testing.T) {
	r := ssehub.NewRoom("a", 8, 8)
	defer r.Close()
	ch, err := r.Subscribe(context.Background(), "s1", 0, 4)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := ssehub.WaitFrame(ctx, ch)
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("want ctx err")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("WaitFrame ignored canceled ctx")
	}
}
