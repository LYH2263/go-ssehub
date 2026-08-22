package ssehub

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRotateAuditReleasesOldFile reproduces the blue-green rotation incident:
// after RotateAudit completes, the previous audit path must be deletable on
// Windows. A leaked handle keeps the process holding the old file open, so
// operators see "in use" errors and the rotation directory piles up.
func TestRotateAuditReleasesOldFile(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old", "audit.log")
	newPath := filepath.Join(dir, "new", "audit.log")

	a, err := openAudit(oldPath)
	if err != nil {
		t.Fatalf("openAudit: %v", err)
	}
	defer a.Close()

	if err := a.Log("first"); err != nil {
		t.Fatalf("Log: %v", err)
	}

	// The rotation target directory already exists in the field (ops create it
	// before rotating); mirror that so this test isolates the handle leak.
	if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		t.Fatalf("mkdir new dir: %v", err)
	}

	if err := a.Rotate(newPath); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	// After rotation the old path must no longer be held open by this process.
	// On Windows a still-open FILE_SHARE_READ handle returns
	// "The process cannot access the file because it is being used by another process."
	if err := os.Remove(oldPath); err != nil {
		t.Fatalf("old path not removable after rotate (leaked handle): %v", err)
	}

	// New path stays writable.
	if err := a.Log("after"); err != nil {
		t.Fatalf("Log after rotate: %v", err)
	}
}

// TestRotateAuditClosesOnPartialFailure ensures that if opening the new file
// fails, the existing file is left untouched and still usable (no half-closed
// state that loses subsequent audit writes).
func TestRotateAuditKeepsOldOnOpenFailure(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "audit.log")

	a, err := openAudit(oldPath)
	if err != nil {
		t.Fatalf("openAudit: %v", err)
	}
	defer a.Close()

	if err := a.Log("first"); err != nil {
		t.Fatalf("Log: %v", err)
	}

	// Target a path whose directory cannot be created (parent is a file).
	badDir := filepath.Join(dir, "blocker")
	if f, err := os.Create(badDir); err != nil {
		t.Fatalf("create blocker: %v", err)
	} else {
		f.Close()
	}
	badPath := filepath.Join(badDir, "child", "audit.log")

	if err := a.Rotate(badPath); err == nil {
		t.Fatalf("Rotate with unwritable target: expected error, got nil")
	}

	// Original handle must still work.
	if err := a.Log("recovered"); err != nil {
		t.Fatalf("Log after failed rotate: %v", err)
	}

	// And still removable once closed.
	a.Close()
	if err := os.Remove(oldPath); err != nil {
		t.Fatalf("old path not removable after failed rotate: %v", err)
	}
}
