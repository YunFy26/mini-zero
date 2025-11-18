package logx

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type mockWriter struct {
	lock    sync.Mutex
	builder strings.Builder
}

func (mw *mockWriter) Alert(v any) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelAlert, v)
}

func (mw *mockWriter) Debug(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelDebug, v, fields...)
}

func (mw *mockWriter) Error(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelError, v, fields...)
}

func (mw *mockWriter) Info(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelInfo, v, fields...)
}

func (mw *mockWriter) Severe(v any) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelSevere, v)
}

func (mw *mockWriter) Slow(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelSlow, v, fields...)
}

func (mw *mockWriter) Stack(v any) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelError, v)
}

func (mw *mockWriter) Stat(v any, fields ...LogField) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	output(&mw.builder, levelStat, v, fields...)
}

func (mw *mockWriter) Close() error {
	return nil
}

func (mw *mockWriter) Contains(text string) bool {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	return strings.Contains(mw.builder.String(), text)
}

func (mw *mockWriter) Reset() {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	mw.builder.Reset()
}

func (mw *mockWriter) String() string {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	return mw.builder.String()
}

func TestAddWriter(t *testing.T) {
	originalLevel := atomic.LoadUint32(&logLevel)
	defer atomic.StoreUint32(&logLevel, originalLevel)
	atomic.StoreUint32(&logLevel, InfoLevel)
	t.Run("SetWriter and Reset", func(t *testing.T) {
		w := new(mockWriter)
		fmt.Printf("1. 创建新的 mockWriter 地址: %p\n", w)

		// 先保存原始的全局写入器（通过 Reset 获取）
		originalWriter := Reset() // 第一次调用 Reset 获取当前全局写入器
		fmt.Printf("2. 原始全局写入器地址: %p\n", originalWriter)

		// 设置新的写入器
		SetWriter(w)
		fmt.Printf("3. 设置 mockWriter 地址: %p\n", w)

		// 重置全局写入器，并获取返回的旧写入器
		oldWriter := Reset()
		fmt.Printf("4. Reset() 返回的旧写入器地址: %p\n", oldWriter)
		fmt.Printf("5. 原来的 w 地址: %p\n", w)

		// 验证
		fmt.Printf("6. Reset() 返回的是否是 w: %v\n", oldWriter == w)
		fmt.Printf("7. 两个地址是否相同: %v\n", fmt.Sprintf("%p", oldWriter) == fmt.Sprintf("%p", w))

		if oldWriter != w {
			t.Errorf("Reset() 应该返回之前设置的写入器，期望: %p, 实际: %p", w, oldWriter)
		}
	})
}

func TestField(t *testing.T) {
	tests := []struct {
		name string
		f    LogField
		want map[string]any
	}{
		{
			name: "error",
			f:    Field("foo", errors.New("bar")),
			want: map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "errors",
			f:    Field("foo", []error{errors.New("bar"), errors.New("baz")}),
			want: map[string]any{
				"foo": []any{"bar", "baz"},
			},
		},
		{
			name: "strings",
			f:    Field("foo", []string{"bar", "baz"}),
			want: map[string]any{
				"foo": []any{"bar", "baz"},
			},
		},
		{
			name: "duration",
			f:    Field("foo", time.Second),
			want: map[string]any{
				"foo": "1s",
			},
		},
		{
			name: "durations",
			f:    Field("foo", []time.Duration{time.Second, 2 * time.Second}),
			want: map[string]any{
				"foo": []any{"1s", "2s"},
			},
		},
		{
			name: "times",
			f: Field("foo", []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			}),
			want: map[string]any{
				"foo": []any{"2020-01-01 00:00:00 +0000 UTC", "2020-01-02 00:00:00 +0000 UTC"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			w := new(mockWriter)
			old := writer.Swap(w)
			defer writer.Store(old)

		})
	}
}
