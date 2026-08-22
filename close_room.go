package ssehub
func (r *Room) Close() {
        r.mu.Lock()
        defer r.mu.Unlock()
        if r.closed { return }
        // flush replay observability before clearing rooms
        _ = r.Replay.Flush()
        for id, s := range r.Subs {
                s.Q.Close()
                delete(r.Subs, id)
        }
        r.Replay.Clear()
        r.closed = true
}
