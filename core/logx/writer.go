package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
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
		// content 脱敏
		val = v.MaskSensitive()
	}
	// 创建日志条目 +3: level, timestamp and content
	entry := make(logEntry, len(fields)+3)
	for _, field := range fields {
		// field value 脱敏
		mval := maskSensitive(field.Value)
		// 对field value进行编码
		entry[field.Key] = processFieldValue(mval)
	}

	// 根据日志格式输出
	switch atomic.LoadUint32(&encoding) {
	case plainEncodingType:
		// 处理key-value结构
		plainFields := buildPlainFields(entry)
		writePlainAny(writer, level, val, plainFields...)
	default:
		entry[timestampKey] = getTimestamp()
		entry[levelKey] = level
		entry[contentKey] = val
		writeJson(writer, entry)
	}
}

// 处理字段值，按照不同类型进行编码
func processFieldValue(value any) any {
	switch val := value.(type) {
	case error:
		return encodeError(val)
	case []error:
		var errs []string
		for _, err := range val {
			errs = append(errs, encodeError(err))
		}
		return errs
	case time.Duration:
		return fmt.Sprint(val)
	case []time.Duration:
		var durs []string
		for _, dur := range val {
			durs = append(durs, fmt.Sprint(dur))
		}
		return durs
	case []time.Time:
		var times []string
		for _, t := range val {
			times = append(times, fmt.Sprint(t))
		}
		return times
	case json.Marshaler:
		return val
	case fmt.Stringer:
		return encodeStringer(val)
	case []fmt.Stringer:
		var strs []string
		for _, s := range val {
			strs = append(strs, encodeStringer(s))
		}
		return strs
	default:
		return val
	}

}

// 把key-value结构展开，构建为字符串切片
func buildPlainFields(fields logEntry) []string {
	items := make([]string, 0, len(fields))
	for k, v := range fields {
		// %v:
		// {Bob 30 {New York NY}}
		// %+v:
		// {Name:Bob Age:30 Address:{City:New York State:NY}}
		items = append(items, fmt.Sprintf("%s=%+v", k, v))
	}
	return items
}

// 写入纯文本格式的日志
func writePlainAny(writer io.Writer, level string, val any, fields ...string) {
	level = wrapLevelWithColor(level)
	switch v := val.(type) {
	case string:
		writePlainText(writer, level, v, fields...)
	case error:
		writePlainText(writer, level, v.Error(), fields...)
	case fmt.Stringer:
		writePlainText(writer, level, v.String(), fields...)
	default:
		writePlainValue(writer, level, val, fields...)
	}
}

func wrapLevelWithColor(level string) string {
	// var colour

	return level
}

// 写入文本日志
func writePlainText(writer io.Writer, level string, msg string, fields ...string) {
	var buf bytes.Buffer
	buf.WriteString(getTimestamp())
	buf.WriteByte(plainEncodingSep)
	buf.WriteString(level)
	buf.WriteByte(plainEncodingSep)
	buf.WriteString(msg)
	for _, field := range fields {
		buf.WriteByte(plainEncodingSep)
		buf.WriteString(field)
	}
	buf.WriteByte('\n')
	if writer == nil {
		log.Println(buf.String())
		return
	}
	if _, err := writer.Write(buf.Bytes()); err != nil {
		log.Println("failed to write log:", err)
	}
}

// 写入key-value结构日志
func writePlainValue(writer io.Writer, level string, val any, fields ...string) {
	var buf bytes.Buffer
	buf.WriteString(getTimestamp())
	buf.WriteByte(plainEncodingSep)
	buf.WriteString(level)
	buf.WriteByte(plainEncodingSep)
	// 将val编码为JSON格式并写入缓冲区
	if err := json.NewEncoder(&buf).Encode(val); err != nil {
		log.Printf("err: %s\n\n%s", err.Error(), debug.Stack())
		return
	}
	for _, field := range fields {
		buf.WriteByte(plainEncodingSep)
		buf.WriteString(field)
	}

	buf.WriteByte('\n')
	if writer == nil {
		log.Println(buf.String())
		return
	}
	if _, err := writer.Write(buf.Bytes()); err != nil {
		log.Println("failed to write log:", err)
	}
}

func writeJson(writer io.Writer, entry logEntry) {

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
