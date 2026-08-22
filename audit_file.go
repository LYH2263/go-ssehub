package ssehub

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type fileAudit struct {
	mu   sync.Mutex
	path string
	f    *os.File
}

func openAudit(path string) (*fileAudit, error) {
	if d := filepath.Dir(path); d != "." {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &fileAudit{path: path, f: f}, nil
}
func (a *fileAudit) Log(line string) error {
	if a == nil || a.f == nil {
		return fmt.Errorf("audit closed")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	_, err := fmt.Fprintln(a.f, line)
	return err
}
func (a *fileAudit) Rotate(newPath string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	f, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if a.f != nil {
		_ = a.f.Close()
	}
	a.f = f
	a.path = newPath
	return nil
}
func (a *fileAudit) Close() error {
	if a == nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.f == nil {
		return nil
	}
	err := a.f.Close()
	a.f = nil
	return err
}
func (h *Hub) EnableAudit(path string) error {
	a, err := openAudit(path)
	if err != nil {
		return err
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.audit = a
	h.auditPath = path
	return nil
}
func (h *Hub) AuditLog(line string) error {
	h.mu.Lock()
	a := h.audit
	h.mu.Unlock()
	if a == nil {
		return nil
	}
	return a.Log(line)
}
func (h *Hub) RotateAudit(path string) error {
	h.mu.Lock()
	a := h.audit
	h.mu.Unlock()
	if a == nil {
		return fmt.Errorf("%w: no audit", ErrInvalid)
	}
	return a.Rotate(path)
}
