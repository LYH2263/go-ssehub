package ssehub
import (
        "context"
        "sync"
        "github.com/LYH2263/go-ssehub/internal/event"
)
type Hub struct {
        mu sync.Mutex
        rooms map[string]*Room
        closed bool
        auditPath string
        audit *fileAudit
}
func NewHub() *Hub { return &Hub{rooms: make(map[string]*Room)} }
func (h *Hub) OpenRoom(name string) *Room {
        h.mu.Lock(); defer h.mu.Unlock()
        if r, ok := h.rooms[name]; ok { return r }
        r := NewRoom(name, 64, 8)
        h.rooms[name] = r
        return r
}
func (h *Hub) Publish(ctx context.Context, room, name string, data []byte) (event.Event, error) {
        return h.OpenRoom(room).Publish(ctx, name, data)
}
func (h *Hub) Close() {
        h.mu.Lock(); defer h.mu.Unlock()
        if h.closed { return }
        h.closed = true
        for _, r := range h.rooms { r.Close() }
        if h.audit != nil { _ = h.audit.Close() }
}
