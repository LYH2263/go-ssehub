package replay

import (
	"bytes"
	"testing"

	"github.com/LYH2263/go-ssehub/internal/event"
)

// TestAppend_DataIsolated reproduces the live-room gateway aliasing bug: the
// caller reuses one payload buffer to push several frames without allocating
// a new slice each time. Before the fix, Append stored the caller's slice
// header verbatim, so rewriting the buffer for the next frame retroactively
// corrupted the earlier frames already sitting in the ring.
func TestAppend_DataIsolated(t *testing.T) {
	r := New(8)

	// One buffer, rewritten between frames — mirrors the gateway's
	// "push a few frames off the same payload to save an allocation" path.
	buf := make([]byte, 5)
	for i := range buf {
		buf[i] = 'a'
	}
	e1 := r.Append(event.Event{Name: "n", Data: buf})
	if e1.ID == 0 {
		t.Fatalf("Append did not assign ID: %+v", e1)
	}

	for i := range buf {
		buf[i] = 'b'
	}
	e2 := r.Append(event.Event{Name: "n", Data: buf})
	if e2.ID == 0 {
		t.Fatalf("Append did not assign ID: %+v", e2)
	}

	// What the ring actually holds must reflect the value at append time,
	// not the buffer's final contents.
	got := r.CloneAll()
	if len(got) != 2 {
		t.Fatalf("CloneAll len = %d, want 2", len(got))
	}
	if want := []byte("aaaaa"); !bytes.Equal(got[0].Data, want) {
		t.Errorf("frame 0 Data = %q, want %q (ring aliased caller buffer)", got[0].Data, want)
	}
	if want := []byte("bbbbb"); !bytes.Equal(got[1].Data, want) {
		t.Errorf("frame 1 Data = %q, want %q", got[1].Data, want)
	}
}

// TestAfter_DataIsolated ensures reconnect replay via Last-Event-ID does not
// hand back data that mutates with the caller's buffer or the ring's later
// appends.
func TestAfter_DataIsolated(t *testing.T) {
	r := New(8)
	id1 := r.Append(event.Event{Name: "n", Data: []byte("one")}).ID
	id2 := r.Append(event.Event{Name: "n", Data: []byte("two")}).ID
	if id1 == id2 {
		t.Fatalf("ids not advancing: %d", id1)
	}

	got := r.After(id1)
	if len(got) != 1 {
		t.Fatalf("After(%d) len = %d, want 1", id1, len(got))
	}
	if want := []byte("two"); !bytes.Equal(got[0].Data, want) {
		t.Errorf("After Data = %q, want %q", got[0].Data, want)
	}
	// Mutating the returned copy must not affect the ring.
	got[0].Data[0] = 'X'
	again := r.After(id1)
	if want := []byte("two"); !bytes.Equal(again[0].Data, want) {
		t.Errorf("ring mutated via After return: got %q, want %q", again[0].Data, want)
	}
}

// TestCloneAll_DataIsolated ensures the snapshot path returns independent
// copies (the old AliasAll returned the internal slice verbatim).
func TestCloneAll_DataIsolated(t *testing.T) {
	r := New(8)
	r.Append(event.Event{Name: "n", Data: []byte("first")})
	r.Append(event.Event{Name: "n", Data: []byte("second")})

	snap := r.CloneAll()
	if len(snap) != 2 {
		t.Fatalf("CloneAll len = %d, want 2", len(snap))
	}
	// Mutate the snapshot; the ring must be unaffected.
	snap[0].Data[0] = 'X'
	again := r.CloneAll()
	if want := []byte("first"); !bytes.Equal(again[0].Data, want) {
		t.Errorf("ring mutated via CloneAll return: got %q, want %q", again[0].Data, want)
	}
}
