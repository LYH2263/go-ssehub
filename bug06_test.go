package ssehub_test

import (
	"context"
	ssehub "github.com/LYH2263/go-ssehub"
	"path/filepath"
	"testing"
)

func TestBug06_AuditFailNoAck(t *testing.T) {
	h := ssehub.NewHub()
	defer h.Close()
	p := filepath.Join(t.TempDir(), "a.log")
	if err := h.EnableAudit(p); err != nil {
		t.Fatal(err)
	}
	if err := h.CloseAuditFile(); err != nil {
		t.Fatal(err)
	}
	if _, err := h.PublishAck(context.Background(), "r1", "e", []byte("x")); err == nil {
		t.Fatal("want fail")
	}
}
