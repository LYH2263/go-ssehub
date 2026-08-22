package ssehub_test

import (
	ssehub "github.com/LYH2263/go-ssehub"
	"os"
	"path/filepath"
	"testing"
)

func TestBug09_RotateClosesOldFile(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.log")
	p2 := filepath.Join(dir, "b.log")
	h := ssehub.NewHub()
	defer h.Close()
	if err := h.EnableAudit(p1); err != nil {
		t.Fatal(err)
	}
	if err := h.AuditLog("one"); err != nil {
		t.Fatal(err)
	}
	if err := h.RotateAudit(p2); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(p1); err != nil {
		t.Fatalf("locked %v", err)
	}
}
