package ssehub
import (
        "sync"
        "github.com/LYH2263/go-ssehub/internal/backpressure"
        "github.com/LYH2263/go-ssehub/internal/event"
        "github.com/LYH2263/go-ssehub/internal/replay"
)
type Sub struct {
        ID string
        Room string
        Q *backpressure.Queue
}
type Room struct {
        mu sync.Mutex
        Name string
        Subs map[string]*Sub
        Replay *replay.Ring
        Enc event.Encoder
        closed bool
}
func NewRoom(name string, replayCap, qCap int) *Room {
        return &Room{
                Name: name,
                Subs: make(map[string]*Sub),
                Replay: replay.New(replayCap),
                Enc: event.DefaultEncoder{},
        }
}
func (r *Room) SetEncoder(enc event.Encoder) {
        r.mu.Lock(); defer r.mu.Unlock()
        r.Enc = enc
}
func (r *Room) ListSubscribers() []string {
        r.mu.Lock(); defer r.mu.Unlock()
        out := make([]string, 0, len(r.Subs))
        for id := range r.Subs { out = append(out, id) }
        return out
}
