package ssehub

import (
	"context"
	"fmt"

	"github.com/LYH2263/go-ssehub/internal/event"
)

func (h *Hub) Fanout(ctx context.Context, room, name string, data []byte) (event.Event, error) {
	ev, err := h.OpenRoom(room).Publish(ctx, name, data)
	if err != nil {
		// %w preserves the error chain so callers can match the
		// backpressure sentinel via errors.Is(err, backpressure.ErrFull)
		// to drive rate limiting / circuit breaking on a full queue.
		return ev, fmt.Errorf("fanout: %w", err)
	}
	return ev, nil
}
