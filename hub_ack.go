package ssehub

import (
	"context"
	"fmt"

	"github.com/LYH2263/go-ssehub/internal/event"
)

func (h *Hub) PublishAck(ctx context.Context, room, name string, data []byte) (event.Event, error) {
	ev, err := h.Publish(ctx, room, name, data)
	if err != nil {
		return ev, err
	}
	_ = h.AuditLog(fmt.Sprintf("ack %s %d", room, ev.ID))
	return ev, nil
}

// CloseAuditFile closes the audit file handle but keeps audit enabled (Log will fail).
func (h *Hub) CloseAuditFile() error {
	h.mu.Lock()
	a := h.audit
	h.mu.Unlock()
	if a == nil {
		return nil
	}
	return a.Close()
}
