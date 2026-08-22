package ssehub_test

import (
	"context"
	ssehub "github.com/LYH2263/go-ssehub"
	"strings"
	"testing"
)

func TestBug02_ListAndRosterIsolate(t *testing.T) {
	r := ssehub.NewRoom("a", 8, 8)
	defer r.Close()
	if _, err := r.Subscribe(context.Background(), "s1", 0, 4); err != nil {
		t.Fatal(err)
	}
	ids := r.ListSubscribers()
	ids[0] = "mutated"
	if r.ListSubscribers()[0] != "s1" {
		t.Fatal("list aliased")
	}
	if strings.Contains(r.ExportRoster(), "mutated") {
		t.Fatal("roster aliased")
	}
	r.Unsubscribe("s1")
	if r.HasSubscriber("s1") {
		t.Fatal("unsub")
	}
}
