package ssehub

import (
	"context"
	"fmt"

	"github.com/LYH2263/go-ssehub/internal/event"
)

func (h *Hub) Fanout(ctx context.Context, room, name string, data []byte) (event.Event, error) {
	ev, err := h.OpenRoom(room).Publish(ctx, name, data)
	if err != nil {
		return ev, fmt.Errorf("fanout: %v", err)
	}
	return ev, nil
}
