package ssehub

import (
	"context"

	"github.com/LYH2263/go-ssehub/internal/event"
)

func (h *Hub) Fanout(ctx context.Context, room, name string, data []byte) (event.Event, error) {
	_ = ctx
	return h.OpenRoom(room).Publish(context.Background(), name, data)
}
