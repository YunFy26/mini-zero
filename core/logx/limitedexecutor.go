package logx

import (
	"sync/atomic"
	"time"

	"github.com/YunFy26/mini-zero/core/syncx"
	"github.com/YunFy26/mini-zero/core/timex"
)

// 限制执行器(在一定时间内只执行一次)
// 控制日志输出频率
type limitedExecutor struct {
	threshold time.Duration
	lastTime  *syncx.AtomicDuration
	discarded uint32 // 记录被丢弃的操作次数
}

func newLimitedExecutor(milliseconds int) *limitedExecutor {
	return &limitedExecutor{
		threshold: time.Duration(milliseconds) * time.Millisecond,
		lastTime:  syncx.NewAtomicDuration(),
	}
}

func (le *limitedExecutor) logOrDiscard(execute func()) {
	if le == nil || le.threshold <= 0 {
		execute()
		return
	}
	now := timex.Now()
	if now-le.lastTime.Load() <= le.threshold {
		atomic.AddUint32(&le.discarded, 1)
	} else {
		le.lastTime.Set(now)
		discarded := atomic.SwapUint32(&le.discarded, 0)
		if discarded > 0 {
			// Errorf("Discarded %d error messages", discarded)
		}
		execute()
	}
}
