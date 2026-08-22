package ssehub

import (
	"context"
	"errors"
	"testing"
	"time"
)

// 带截止时间的等帧在 ctx 超时后必须立即收工，不能死堵在接收通道上耗着。
// 回归用例：超时后仍在等下一帧会把整条网关链路的超时拖长好几秒。
func TestWaitFrame_ContextDeadlineReturnsImmediately(t *testing.T) {
	ch := make(chan []byte) // 永不主动发帧，模拟“等下一帧”的死堵场景

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	b, err := WaitFrame(ctx, ch)
	elapsed := time.Since(start)

	if b != nil {
		t.Fatalf("expected nil frame on timeout, got %v", b)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	// 必须在截止时间附近收工，而不是被接收通道拖成秒级。
	if elapsed > 500*time.Millisecond {
		t.Fatalf("WaitFrame did not return shortly after deadline: %v", elapsed)
	}
}

// 主动取消同样应当尽快返回，不能感觉不到取消。
func TestWaitFrame_ContextCancelReturnsImmediately(t *testing.T) {
	ch := make(chan []byte)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	b, err := WaitFrame(ctx, ch)
	elapsed := time.Since(start)

	if b != nil {
		t.Fatalf("expected nil frame on cancel, got %v", b)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected Canceled, got %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("WaitFrame did not return shortly after cancel: %v", elapsed)
	}
}

// 收到帧时仍应正常返回，不能因为引入 ctx 监听而影响正常路径。
func TestWaitFrame_ReturnsFrame(t *testing.T) {
	ch := make(chan []byte, 1)
	ch <- []byte("frame")

	b, err := WaitFrame(context.Background(), ch)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if string(b) != "frame" {
		t.Fatalf("expected frame, got %q", b)
	}
}

// 通道关闭应返回 ErrClosed。
func TestWaitFrame_ClosedChannel(t *testing.T) {
	ch := make(chan []byte)
	close(ch)

	b, err := WaitFrame(context.Background(), ch)
	if b != nil {
		t.Fatalf("expected nil frame on closed channel, got %v", b)
	}
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("expected ErrClosed, got %v", err)
	}
}
