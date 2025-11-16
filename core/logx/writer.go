package logx

import (
	"io"
	"log"
	"sync"
	"sync/atomic"
)

type (
	Writer interface {
		Alert(v any)                     // 警告日志
		Close() error                    // 关闭写入器
		Debug(v any, fields ...LogField) // 调试日志
		Error(v any, fields ...LogField) // 错误日志
		Info(v any, fields ...LogField)  // 信息日志
		Severe(v any)                    // 严重日志
		Slow(v any, fields ...LogField)  // 慢查询日志
		Stack(v any)                     // 堆栈日志
		Stat(v any, fields ...LogField)  // 统计日志
	}

	// 原子写入
	atomicWriter struct {
		writer Writer
		lock   sync.RWMutex
	}
	// 组合写入
	comboWriter struct {
		writer []Writer
	}
	// 具体的日志写入器实现
	concreteWriter struct {
		infoLog   io.WriteCloser
		errorLog  io.WriteCloser
		severeLog io.WriteCloser
		slowLog   io.WriteCloser
		statLog   io.WriteCloser
		stackLog  io.Writer
	}

	nopWriter struct{}
)

func NewWriter(w io.Writer) Writer {
	lw := newLogWriter(log.New(w, "", flags))

	return &concreteWriter{
		infoLog:   lw,
		errorLog:  lw,
		severeLog: lw,
		slowLog:   lw,
		statLog:   lw,
		stackLog:  lw,
	}
}

func (w *concreteWriter) Alert(v any) {
	output(w.errorLog, levelAlert, v)
}

func (w *concreteWriter) Close() error {
	return nil
}

func (w *concreteWriter) Debug(v any, fields ...LogField) {

}

func (w *concreteWriter) Error(v any, fields ...LogField) {

}

func (w *concreteWriter) Info(v any, fields ...LogField) {

}

func (w *concreteWriter) Severe(v any) {

}

func (w *concreteWriter) Slow(v any, fields ...LogField) {

}

func (w *concreteWriter) Stack(v any) {

}

func (w *concreteWriter) Stat(v any, fields ...LogField) {

}

func output(writer io.Writer, level string, val any, fields ...LogField) {
	switch v := val.(type) {
	case string:
		// 检查是否需要截断
		maxLen := atomic.LoadUint32(&maxContentLength)
		if maxLen > 0 && len(v) > int(maxLen) {
			val = v[:maxLen]
			fields = append(fields, truncatedField)
		}

	case Sensitive:
		// 脱敏
		val = v.MaskSensitive()
	}
	// 创建日志条目
	// entry := make(logEntry, len(fields)+3)
	// for _, field := range fields {
	// 	mval := maskSensitive(field.Value)
	// 	// entry[field.Key] = processFieldValue(mval)
	// }

	// 根据编码格式输出

}

func (n nopWriter) Alert(_ any)                {}
func (n nopWriter) Close() error               { return nil }
func (n nopWriter) Debug(_ any, _ ...LogField) {}
func (n nopWriter) Error(_ any, _ ...LogField) {}
func (n nopWriter) Info(_ any, _ ...LogField)  {}
func (n nopWriter) Severe(_ any)               {}
func (n nopWriter) Slow(_ any, _ ...LogField)  {}
func (n nopWriter) Stack(_ any)                {}
func (n nopWriter) Stat(_ any, _ ...LogField)  {}
