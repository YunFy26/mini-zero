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

	"github.com/YunFy26/mini-zero/core/color"
	"github.com/YunFy26/mini-zero/core/errorx"
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
		writers []Writer
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

func (w *atomicWriter) Load() Writer {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.writer
}

func (w *atomicWriter) Store(writer Writer) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.writer = writer
}

func (w *atomicWriter) StoreIfNil(v Writer) Writer {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.writer == nil {
		w.writer = v
	}
	return w.writer
}

// 切换写入器，返回旧的写入器
func (w *atomicWriter) Swap(v Writer) Writer {
	w.lock.Lock()
	defer w.lock.Unlock()
	old := w.writer
	w.writer = v
	return old
}

func (c comboWriter) Alert(v any) {
	for _, w := range c.writers {
		w.Alert(v)
	}
}

func (c comboWriter) Close() error {
	var be errorx.BatchError
	for _, w := range c.writers {
		be.Add(w.Close())
	}
	return be.Err()
}

func (c comboWriter) Debug(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Debug(v, fields...)
	}
}

func (c comboWriter) Error(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Error(v, fields...)
	}
}

func (c comboWriter) Info(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Info(v, fields...)
	}
}

func (c comboWriter) Severe(v any) {
	for _, w := range c.writers {
		w.Severe(v)
	}
}

func (c comboWriter) Slow(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Slow(v, fields...)
	}
}

func (c comboWriter) Stack(v any) {
	for _, w := range c.writers {
		w.Stack(v)
	}
}

func (c comboWriter) Stat(v any, fields ...LogField) {
	for _, w := range c.writers {
		w.Stat(v, fields...)
	}
}

// func newConsoleWriter() Writer {
// 	outLog := newLogWriter(log.New(fatihcolor.Output, "", flags))
// 	errLog := newLogWriter(log.New(fatihcolor.Error, "", flags))
// 	return &concreteWriter{
// 		infoLog:   outLog,
// 		errorLog:  errLog,
// 		severeLog: errLog,
// 		slowLog:   outLog,
// 		statLog:   outLog,
// 		stackLog:  newLessWriter(errLog, options.logStackCooldownMills),
// 	}
// }

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
		// 格式化输出 %v 和 %+v 的区别
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
	var colour color.Color
	switch level {
	case levelAlert:
		colour = color.FgRed
	case levelError:
		colour = color.FgRed
	case levelFatal:
		colour = color.FgRed
	case levelInfo:
		colour = color.FgGreen
	case levelSlow:
		colour = color.FgYellow
	case levelDebug:
		colour = color.FgYellow
	case levelStat:
		colour = color.FgGreen
	}

	if colour == color.NoColor {
		return level
	}

	return color.WithColorPadding(level, colour)
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

func writeJson(writer io.Writer, info any) {
	if content, err := marshalJson(info); err != nil {
		log.Printf("err: %s\n\n%s", err.Error(), debug.Stack())
	} else if writer == nil {
		log.Println(string(contentKey))
	} else {
		if _, err := writer.Write(append(content, '\n')); err != nil {
			log.Println(err.Error())
		}
	}
}

// 将任意值编码为JSON格式的字节切片
func marshalJson(v any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	// 禁止HTML转义
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	// 移除末尾的换行符
	if l := buf.Len(); l > 0 && buf.Bytes()[l-1] == '\n' {
		buf.Truncate(l - 1)
	}
	return buf.Bytes(), err
}

// TODO: 实现全局字段合并
func mergeGloablFields(fields []LogField) []LogField {
	return nil
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
