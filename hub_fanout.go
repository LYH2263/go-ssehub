package ssehub

import (
	"context"

	"github.com/LYH2263/go-ssehub/internal/event"
)

func (h *Hub) Fanout(ctx context.Context, room, name string, data []byte) (event.Event, error) {
	if err := ctx.Err(); err != nil {
		return event.Event{}, err
	}
	return h.OpenRoom(room).Publish(ctx, name, data)
}
