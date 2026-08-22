package ssehub_test

import (
	"context"
	"testing"

	ssehub "github.com/LYH2263/go-ssehub"
)

func TestBug02_ListSubscribersIsolated(t *testing.T) {
	r := ssehub.NewRoom("a", 8, 8)
	defer r.Close()
	if _, err := r.Subscribe(context.Background(), "s1", 0, 4); err != nil {
		t.Fatal(err)
	}
	ids := r.ListSubscribers()
	ids[0] = "mutated"
	if r.ListSubscribers()[0] != "s1" {
		t.Fatal("aliased")
	}
}
