package timex

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	// 测试 Now() 返回的时间应该是相对于 initTime 的持续时间
	now1 := Now()

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	now2 := Now()

	// 第二次调用应该比第一次大
	if now2 <= now1 {
		t.Fatalf("Expected now2 > now1, got now1=%v, now2=%v", now1, now2)
	}

	// 时间差应该大致等于我们等待的时间（允许一些误差）
	diff := now2 - now1
	if diff < 5*time.Millisecond || diff > 20*time.Millisecond {
		t.Fatalf("Expected diff around 10ms, got %v", diff)
	}
}

func TestSince(t *testing.T) {
	start := Now()

	// 等待一段时间
	time.Sleep(15 * time.Millisecond)

	// Since 应该返回从 start 到现在经过的时间
	elapsed := Since(start)

	// 经过的时间应该大致等于等待时间
	if elapsed < 10*time.Millisecond || elapsed > 25*time.Millisecond {
		t.Fatalf("Expected elapsed around 15ms, got %v", elapsed)
	}

	// 验证 Now() - start 应该等于 Since(start)
	now := Now()
	if now-start != elapsed {
		t.Fatalf("Expected Now()-start == Since(start), got %v != %v", now-start, elapsed)
	}
}
