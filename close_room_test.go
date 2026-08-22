package ssehub

import (
	"context"
	"testing"
)

// TestCloseFlushCount_ReflectsRealPendingReplay 回归末班对账：
// 房间重放环里压着若干帧时关闭，CloseFlushCount 必须返回真实刷出量，
// 不能被随后的清空动作抹成 0。
func TestCloseFlushCount_ReflectsRealPendingReplay(t *testing.T) {
	r := NewRoom("night-shift", 64, 8)

	// 往重放环里压入 3 帧，不留订阅者（末班关房前的典型状态）。
	for i := 0; i < 3; i++ {
		if _, err := r.Publish(context.Background(), "tick", []byte("frame")); err != nil {
			t.Fatalf("publish %d: %v", i, err)
		}
	}
	if got := r.PendingReplay(); got != 3 {
		t.Fatalf("before close, pending replay = %d, want 3", got)
	}

	got := r.CloseFlushCount()
	const want = 3
	if got != want {
		t.Fatalf("CloseFlushCount = %d, want %d (flushed-before-clear order matters)", got, want)
	}

	// 关房后再次调用应返回 0（房间已 closed，幂等短路）。
	if second := r.CloseFlushCount(); second != 0 {
		t.Fatalf("second CloseFlushCount = %d, want 0 after close", second)
	}
}

// TestCloseFlushCount_EmptyRoom 零帧关房计数为 0，不误报。
func TestCloseFlushCount_EmptyRoom(t *testing.T) {
	r := NewRoom("idle", 64, 8)
	if got := r.CloseFlushCount(); got != 0 {
		t.Fatalf("empty room CloseFlushCount = %d, want 0", got)
	}
}
