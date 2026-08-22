package ssehub
import ("fmt"; "sync")
type Hub struct {
        mu sync.Mutex
        subs map[string]chan []byte
        closed bool
}
func NewHub() *Hub { return &Hub{subs: make(map[string]chan []byte)} }
func (h *Hub) Subscribe(id string, buf int) (<-chan []byte, error) {
        if buf <= 0 { buf = 8 }
        h.mu.Lock(); defer h.mu.Unlock()
        if h.closed { return nil, ErrClosed }
        if _, ok := h.subs[id]; ok { return nil, ErrConflict }
        ch := make(chan []byte, buf)
        h.subs[id] = ch
        return ch, nil
}
func (h *Hub) Publish(id string, payload []byte) error {
        h.mu.Lock(); defer h.mu.Unlock()
        if h.closed { return ErrClosed }
        ch, ok := h.subs[id]
        if !ok { return ErrNotFound }
        select {
        case ch <- append([]byte(nil), payload...):
                return nil
        default:
                return fmt.Errorf("%w: backlog full", ErrInvalid)
        }
}
func (h *Hub) Unsubscribe(id string) {
        h.mu.Lock(); defer h.mu.Unlock()
        if ch, ok := h.subs[id]; ok {
                close(ch); delete(h.subs, id)
        }
}
func (h *Hub) Close() {
        h.mu.Lock(); defer h.mu.Unlock()
        if h.closed { return }
        h.closed = true
        for id, ch := range h.subs { close(ch); delete(h.subs, id) }
}
