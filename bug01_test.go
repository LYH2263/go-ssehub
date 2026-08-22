package ssehub_test

import (
	"bytes"
	"context"
	ssehub "github.com/LYH2263/go-ssehub"
	"testing"
)

func TestBug01_PublishAndSnapshotIsolate(t *testing.T) {
	r := ssehub.NewRoom("a", 8, 8)
	defer r.Close()
	buf := []byte("alpha")
	if _, err := r.Publish(context.Background(), "e", buf); err != nil {
		t.Fatal(err)
	}
	buf[0] = 'A'
	snap := r.SnapshotEvents()
	if len(snap) == 0 || bytes.Equal(snap[0].Data, buf) {
		t.Fatalf("publish leak %q", snap[0].Data)
	}
	snap[0].Data[0] = 'Z'
	if r.SnapshotEvents()[0].Data[0] == 'Z' {
		t.Fatal("snapshot leak")
	}
}
