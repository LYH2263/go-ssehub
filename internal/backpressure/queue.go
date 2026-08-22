package backpressure
import ("errors"; "sync")
var ErrFull = errors.New("backpressure: full")
type Queue struct {
        mu sync.Mutex
        ch chan []byte
        closed bool
}
func New(n int) *Queue {
        if n <= 0 { n = 8 }
        return &Queue{ch: make(chan []byte, n)}
}
func (q *Queue) Push(b []byte) error {
        q.mu.Lock()
        defer q.mu.Unlock()
        if q.closed { return errors.New("backpressure: closed") }
        select {
        case q.ch <- append([]byte(nil), b...):
                return nil
        default:
                return ErrFull
        }
}
func (q *Queue) Chan() <-chan []byte { return q.ch }
func (q *Queue) Close() {
        q.mu.Lock(); defer q.mu.Unlock()
        if q.closed { return }
        q.closed = true
        close(q.ch)
}
