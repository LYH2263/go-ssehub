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
	if err := h.AuditLog(fmt.Sprintf("ack %s %d", room, ev.ID)); err != nil {
		return event.Event{}, fmt.Errorf("%w: audit", err)
	}
	return ev, nil
}
