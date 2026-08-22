package ssehub

import (
	"github.com/LYH2263/go-ssehub/internal/backpressure"
	"github.com/LYH2263/go-ssehub/internal/event"
	"github.com/LYH2263/go-ssehub/internal/replay"
	"sync"
)

type Sub struct {
	ID   string
	Room string
	Q    *backpressure.Queue
}
type Room struct {
	mu       sync.Mutex
	Name     string
	Subs     map[string]*Sub
	subOrder []string
	Replay   *replay.Ring
	Enc      event.Encoder
	closed   bool
}

func NewRoom(name string, replayCap, qCap int) *Room {
	return &Room{
		Name:   name,
		Subs:   make(map[string]*Sub),
		Replay: replay.New(replayCap),
		Enc:    event.DefaultEncoder{},
	}
}
func (r *Room) SetEncoder(enc event.Encoder) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Enc = enc
}
func (r *Room) ListSubscribers() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.subOrder))
	copy(out, r.subOrder)
	return out
}

func (r *Room) PendingReplay() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Replay == nil {
		return 0
	}
	return r.Replay.Len()
}
