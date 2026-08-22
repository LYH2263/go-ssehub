package ssehub

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// tmpAudit writes a fresh audit file under a temp dir and returns its path.
func tmpAudit(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	p := filepath.Join(d, "audit.log")
	if err := os.WriteFile(p, nil, 0o644); err != nil {
		t.Fatalf("seed audit file: %v", err)
	}
	return p
}

// reproduce the compliance-migration incident: ops closes the audit file
// handle to relocate the directory, PublishAck must NOT report success while
// the audit side is unwritable.
func TestPublishAckFailsWhenAuditHandleClosed(t *testing.T) {
	h := NewHub()
	h.OpenRoom("orders")

	p := tmpAudit(t)
	if err := h.EnableAudit(p); err != nil {
		t.Fatalf("EnableAudit: %v", err)
	}

	// baseline: with the handle open, PublishAck returns the audited event.
	ev, err := h.PublishAck(context.Background(), "orders", "created", []byte("x"))
	if err != nil {
		t.Fatalf("first PublishAck: %v", err)
	}
	if ev.ID == 0 {
		t.Fatal("first PublishAck returned zero-id event")
	}

	// ops closes the handle mid-migration; audit stays "enabled" but unwritable.
	if err := h.CloseAuditFile(); err != nil {
		t.Fatalf("CloseAuditFile: %v", err)
	}

	// second PublishAck must surface the audit failure, not success.
	ev2, err := h.PublishAck(context.Background(), "orders", "created", []byte("y"))
	if !errors.Is(err, ErrAudit) {
		t.Fatalf("PublishAck with closed audit: want ErrAudit, got %v", err)
	}
	if ev2.ID == 0 {
		t.Fatal("PublishAck still returned the event despite audit failure")
	}

	// the audit file must contain only the first (durable) line, proving the
	// second write never landed.
	got, _ := os.ReadFile(p)
	if lineCount(string(got)) != 1 {
		t.Fatalf("audit lines = %d, want 1 (closed-handle PublishAck must not land): %q", lineCount(string(got)), string(got))
	}
}

// RotateAudit reopens a writable handle, so PublishAck recovers after a close.
func TestPublishAckRecoversAfterRotate(t *testing.T) {
	h := NewHub()
	h.OpenRoom("orders")

	p1 := tmpAudit(t)
	if err := h.EnableAudit(p1); err != nil {
		t.Fatalf("EnableAudit: %v", err)
	}
	if _, err := h.PublishAck(context.Background(), "orders", "e", []byte("1")); err != nil {
		t.Fatalf("PublishAck before rotate: %v", err)
	}

	if err := h.CloseAuditFile(); err != nil {
		t.Fatalf("CloseAuditFile: %v", err)
	}

	// still closed -> failure
	if _, err := h.PublishAck(context.Background(), "orders", "e", []byte("2")); !errors.Is(err, ErrAudit) {
		t.Fatalf("PublishAck after close: want ErrAudit, got %v", err)
	}

	// rotate onto a fresh path -> recovered
	p2 := tmpAudit(t)
	if err := h.RotateAudit(p2); err != nil {
		t.Fatalf("RotateAudit: %v", err)
	}
	if _, err := h.PublishAck(context.Background(), "orders", "e", []byte("3")); err != nil {
		t.Fatalf("PublishAck after rotate: %v", err)
	}

	got, _ := os.ReadFile(p2)
	if lineCount(string(got)) != 1 {
		t.Fatalf("rotated audit lines = %d, want 1: %q", lineCount(string(got)), string(got))
	}

	// release the audit handle so TempDir cleanup succeeds on Windows.
	if err := h.CloseAuditFile(); err != nil {
		t.Fatalf("final CloseAuditFile: %v", err)
	}
}

// with no audit configured at all, PublishAck succeeds (no audit contract).
func TestPublishAckWithoutAudit(t *testing.T) {
	h := NewHub()
	h.OpenRoom("orders")
	if _, err := h.PublishAck(context.Background(), "orders", "e", []byte("x")); err != nil {
		t.Fatalf("PublishAck without audit: %v", err)
	}
}

// lineCount counts newline-terminated lines in s.
func lineCount(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			n++
		}
	}
	return n
}
