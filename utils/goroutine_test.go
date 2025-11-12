package utils

import (
	"sync"
	"testing"
)

func TestGetGoroutineID(t *testing.T) {
	// 获取主 goroutine 的 ID
	mainID := GetGoroutineID()
	if mainID <= 0 {
		t.Errorf("GetGoroutineID() returned invalid ID: %d", mainID)
	}
	t.Logf("Main goroutine ID: %d", mainID)

	// 测试在多个 goroutine 中获取 ID
	var wg sync.WaitGroup
	ids := make(map[int64]bool)
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := GetGoroutineID()
			mu.Lock()
			ids[id] = true
			mu.Unlock()
			t.Logf("Goroutine ID: %d", id)
		}()
	}

	wg.Wait()

	// 验证每个 goroutine 都有唯一的 ID
	if len(ids) < 10 {
		t.Errorf("Expected at least 10 unique goroutine IDs, got %d", len(ids))
	}
}

func TestGetGoroutineIDString(t *testing.T) {
	idStr := GetGoroutineIDString()
	if idStr == "" {
		t.Errorf("GetGoroutineIDString() returned empty string")
	}
	t.Logf("Goroutine ID (string): %s", idStr)
}

func TestGetGoroutineIDFast(t *testing.T) {
	mainID := GetGoroutineIDFast()
	if mainID <= 0 {
		t.Errorf("GetGoroutineIDFast() returned invalid ID: %d", mainID)
	}
	t.Logf("Main goroutine ID (fast): %d", mainID)

	// 验证两种方法返回相同的 ID
	id1 := GetGoroutineID()
	id2 := GetGoroutineIDFast()
	if id1 != id2 {
		t.Errorf("GetGoroutineID() = %d, GetGoroutineIDFast() = %d, want same", id1, id2)
	}
}

func BenchmarkGetGoroutineID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGoroutineID()
	}
}

func BenchmarkGetGoroutineIDFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGoroutineIDFast()
	}
}
