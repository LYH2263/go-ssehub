package ssehub_test

import (
	"context"
	"errors"
	ssehub "github.com/LYH2263/go-ssehub"
	"testing"
)

func TestBug03_PublishAfterClose(t *testing.T) {
	r := ssehub.NewRoom("a", 8, 8)
	r.Close()
	_, err := r.Publish(context.Background(), "x", []byte("z"))
	if !errors.Is(err, ssehub.ErrClosed) {
		t.Fatalf("%v", err)
	}
}
