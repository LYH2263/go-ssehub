package ssehub

import (
	"context"
	"strings"
	"testing"
)

// TestListSubscribersCopyIsolatesRosterAndUnsubscribe reproduces the "值班控制台"
// scenario: a caller takes the subscriber list, rewrites one element to a
// display nickname, and the same process then exports the roster and tears
// down the original subscription. The presentation copy must not leak into
// the roster nor break Unsubscribe's lookup by the original ID.
func TestListSubscribersCopyIsolatesRosterAndUnsubscribe(t *testing.T) {
	r := NewRoom("war-room", 8, 8)
	ctx := context.Background()

	const origID = "sub-7"
	if _, err := r.Subscribe(ctx, origID, 0, 8); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	// a couple of neighbours so ordering is non-trivial
	if _, err := r.Subscribe(ctx, "sub-8", 0, 8); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	if _, err := r.Subscribe(ctx, "sub-9", 0, 8); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	// Console hands out a view and rewrites one entry to a display nickname.
	view := r.ListSubscribers()
	idx := indexOf(view, origID)
	if idx < 0 {
		t.Fatalf("origID missing from view: %v", view)
	}
	view[idx] = "值班-昵称"

	// Sanity: the authoritative roster must still contain the real ID, not the alias.
	roster := r.ExportRoster()
	if strings.Contains(roster, "值班-昵称") {
		t.Fatalf("roster leaked display alias from presentation copy: %q", roster)
	}
	if !strings.Contains(roster, origID) {
		t.Fatalf("roster lost original id: %q", roster)
	}

	// Unsubscribe by the ORIGINAL id must still succeed and the id must be gone
	// from the roster afterwards.
	r.Unsubscribe(origID)
	after := r.ExportRoster()
	if strings.Contains(after, origID) {
		t.Fatalf("origID still present after Unsubscribe: %q", after)
	}
	if strings.Contains(after, "值班-昵称") {
		t.Fatalf("alias present in roster after Unsubscribe: %q", after)
	}
	if !r.HasSubscriber(origID) {
		// expected: it's gone
	} else {
		t.Fatalf("HasSubscriber should be false after Unsubscribe")
	}
	// remaining subscribers untouched
	if !r.HasSubscriber("sub-8") || !r.HasSubscriber("sub-9") {
		t.Fatalf("neighbours wrongly dropped: %q", after)
	}
}

// TestListSubscribersMutationDoesNotCorruptSource is a focused guard: mutating
// the returned slice in any position must never alter r.subOrder.
func TestListSubscribersMutationDoesNotCorruptSource(t *testing.T) {
	r := NewRoom("room-x", 8, 8)
	ctx := context.Background()
	for _, id := range []string{"a", "b", "c"} {
		if _, err := r.Subscribe(ctx, id, 0, 8); err != nil {
			t.Fatalf("Subscribe %s: %v", id, err)
		}
	}

	v1 := r.ListSubscribers()
	for i := range v1 {
		v1[i] = "POISON-" + string(rune('A'+i))
	}
	v2 := r.ListSubscribers()
	want := []string{"a", "b", "c"}
	if len(v2) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(v2), len(want), v2)
	}
	for i, got := range v2 {
		if got != want[i] {
			t.Fatalf("v2[%d] = %q, want %q (source corrupted by prior mutation)", i, got, want[i])
		}
	}
}

func indexOf(s []string, v string) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}
