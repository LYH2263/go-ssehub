package ssehub_test

import (
	"context"
	"errors"
	"testing"

	ssehub "github.com/LYH2263/go-ssehub"
)

func TestBug04_SubscribeNilEncoder(t *testing.T) {
	r := ssehub.NewRoom("a", 8, 8)
	defer r.Close()
	r.SetEncoder(nil)
	_, err := r.Subscribe(context.Background(), "s1", 0, 4)
	if err == nil || !errors.Is(err, ssehub.ErrInvalid) {
		t.Fatalf("%v", err)
	}
}
