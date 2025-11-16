package logx

import (
	"errors"

	"github.com/YunFy26/mini-zero/core/syncx"
)

// 日志级别常量
const (
	DebugLevel   uint32 = iota // 0: 调试级别，记录所有日志
	InfoLevel                  // 1: 信息级别，不包含调试信息
	ErrorLevel                 // 2: 错误级别，包含错误、慢查询、堆栈
	SevereLevel                // 3: 严重级别，只记录严重错误
	disableLevel = 0xff        // 255: 禁用所有日志
)

// 编码类型常量
const (
	jsonEncodingType  = iota // 0: JSON 编码
	plainEncodingType        // 1: 纯文本编码
)

// 文件名和模式常量
const (
	// 日志文件名
	accessFilename = "access.log" // 访问日志
	errorFilename  = "error.log"  // 错误日志
	severeFilename = "severe.log" // 严重错误日志
	slowFilename   = "slow.log"   // 慢查询日志
	statFilename   = "stat.log"   // 统计日志

	// 编码方式
	plainEncoding    = "plain" // 纯文本编码
	plainEncodingSep = '\t'    // 纯文本分隔符（制表符）

	// 日志轮转规则
	sizeRotationRule = "size" // 按大小轮转

	// 写入模式
	fileMode   = "file"   // 文件模式
	volumeMode = "volume" // 卷模式（可能指日志卷）

	// 日志级别
	levelAlert  = "alert"  // 警报级别
	levelInfo   = "info"   // 信息级别
	levelError  = "error"  // 错误级别
	levelSevere = "severe" // 严重级别
	levelFatal  = "fatal"  // 致命错误级别
	levelSlow   = "slow"   // 慢查询级别
	levelStat   = "stat"   // 统计级别
	levelDebug  = "debug"  // 调试级别

	backupFileDelimiter = "-"
	nilAngleString      = "<nil>"
	flags               = 0x0
)

// 默认key常量
const (
	defaultCallerKey    = "caller"     // 调用者信息字段名
	defaultContentKey   = "content"    // 日志内容字段名
	defaultDurationKey  = "duration"   // 耗时字段名
	defaultLevelKey     = "level"      // 日志级别字段名
	defaultSpanKey      = "span"       // 跨度ID字段名
	defaultTimestampKey = "@timestamp" // 时间戳字段名
	defaultTraceKey     = "trace"      // 追踪ID字段名
	defaultTruncatedKey = "truncated"  // 截断标记字段名
)

var (
	// 日志路径未设置错误
	ErrLogPathNotSet = errors.New("log path must be set")
	// 日志服务名称未设置错误
	ErrLogServiceNameNotSet = errors.New("log service name must be set")
	// 是否在致命错误时退出
	ExitOnFatal = syncx.ForAtomicBool(true)
	// 标记日志内容是否被截断（日志内容太长时，就会被截断）
	truncatedField = Field(truncatedKey, true)
)

var (
	callerKey    = defaultCallerKey
	contentKey   = defaultContentKey
	durationKey  = defaultDurationKey
	levelKey     = defaultLevelKey
	spanKey      = defaultSpanKey
	timestampKey = defaultTimestampKey
	traceKey     = defaultTraceKey
	truncatedKey = defaultTruncatedKey
)
