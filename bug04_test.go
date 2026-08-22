package ssehub_test

import (
	"context"
	"errors"
	ssehub "github.com/LYH2263/go-ssehub"
	"testing"
)

func TestBug04_NilEncoderNoDirtySub(t *testing.T) {
	r := ssehub.NewRoom("a", 8, 8)
	defer r.Close()
	if _, err := r.Publish(context.Background(), "seed", []byte("s")); err != nil {
		t.Fatal(err)
	}
	r.SetEncoder(nil)
	var panicked bool
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		_, err := r.Subscribe(context.Background(), "s1", 0, 4)
		if panicked {
			return
		}
		if err == nil || !errors.Is(err, ssehub.ErrInvalid) {
			t.Fatalf("%v", err)
		}
	}()
	if panicked {
		t.Fatal("panic")
	}
	if r.HasSubscriber("s1") {
		t.Fatal("dirty sub")
	}
	r.SetEncoder(ssehub.NewDefaultEncoder())
	if _, err := r.Subscribe(context.Background(), "s1", 0, 4); err != nil {
		t.Fatal(err)
	}
}
