package ssehub

import (
	"context"
	"fmt"

	"github.com/LYH2263/go-ssehub/internal/event"
)

// PublishAck publishes an event to the room and then writes an audit record
// for it. The returned error is non-nil only when the write is durably
// acknowledged: a closed/unwritable audit side fails the whole call so the
// upstream caller cannot mistake an unaudited event for an audited one.
func (h *Hub) PublishAck(ctx context.Context, room, name string, data []byte) (event.Event, error) {
	ev, err := h.Publish(ctx, room, name, data)
	if err != nil {
		return ev, err
	}
	if err := h.AuditLog(fmt.Sprintf("ack %s %d", room, ev.ID)); err != nil {
		return ev, fmt.Errorf("%w: %v", ErrAudit, err)
	}
	return ev, nil
}

// CloseAuditFile closes the audit file handle but keeps audit enabled, so a
// subsequent Log will fail. This exists for the compliance-migration path
// where ops must relocate the audit directory. It must NOT be called while
// PublishAck traffic is in flight: after CloseAuditFile, every PublishAck
// returns ErrAudit until EnableAudit/RotateAudit reopens a writable handle.
// Prefer RotateAudit, which swaps the handle atomically without an outage.
func (h *Hub) CloseAuditFile() error {
	h.mu.Lock()
	a := h.audit
	h.mu.Unlock()
	if a == nil {
		return nil
	}
	return a.Close()
}
