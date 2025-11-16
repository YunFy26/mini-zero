package logx

type (
	// LogConf 定义日志系统的配置参数
	//
	// 示例配置：
	//  {
	//      "ServiceName": "user-service",
	//      "Mode": "file",
	//      "Level": "info",
	//      "Path": "/var/logs"
	//  }
	LogConf struct {
		// ServiceName 服务名称，用于标识日志来源
		//
		// 示例：
		//  - "user-service"
		//  - "payment-api"
		//
		// 该字段会作为日志字段输出，便于日志聚合和检索
		ServiceName string `json:",optional"`

		// Mode 日志输出模式
		//
		// 可选值：
		//  - "console": 输出到控制台，适用于开发环境
		//  - "file":    输出到文件，适用于生产环境
		//  - "volume":  在K8s环境中使用，在文件名前添加主机名
		//
		// 默认值: "console"
		Mode string `json:",default=console,options=[console,file,volume]"`

		// Encoding 日志编码格式
		//
		// 可选值：
		//  - "json":  JSON格式，机器可读，适合生产环境
		//  - "plain": 纯文本格式，人类可读，适合开发环境
		//
		// 默认值: "json"
		Encoding string `json:",default=json,options=[json,plain]"`

		// TimeFormat 日志时间格式
		//
		// 遵循Go时间格式规范，使用参考时间定义格式：
		//  Mon Jan 2 15:04:05 MST 2006
		//
		// 示例：
		//  - "2006-01-02 15:04:05"
		//  - "2006-01-02T15:04:05.000Z07:00" (ISO8601)
		//
		// 默认值: "2006-01-02T15:04:05.000Z07:00"
		TimeFormat string `json:",optional"`

		// Path 日志文件存储路径
		//
		// 当Mode为"file"或"volume"时生效
		//
		// 示例：
		//  - "logs"              // 相对路径
		//  - "/var/log/myapp"     // 绝对路径
		//
		// 默认值: "logs"
		Path string `json:",default=logs"`

		// Level 日志级别，用于过滤日志记录
		//
		// 级别顺序（从低到高）：
		//  debug < info < error < severe
		//
		// 设置某个级别后，只会记录该级别及更高级别的日志
		// 示例：设置为"info"时，会记录info、error、severe级别日志
		//
		// 默认值: "info"
		Level string `json:",default=info,options=[debug,info,error,severe]"`

		// MaxContentLength 单条日志内容最大字节数
		//
		// 用于防止日志过大，超过此长度的日志内容会被截断
		// 设置为0表示无限制
		//
		// 示例：
		//  - 1024: 限制单条日志最多1024字节
		//  - 0:    无限制
		MaxContentLength uint32 `json:",optional"`

		// Compress 是否压缩旧的日志文件
		//
		// 开启后可节省磁盘空间，压缩格式为gzip
		//
		// 示例：
		//  - true:  启用压缩
		//  - false: 禁用压缩（默认）
		Compress bool `json:",optional"`

		// Stat 是否记录统计日志
		//
		// 统计日志包括QPS、耗时等系统指标
		//
		// 生产环境建议开启，开发环境可关闭以减少日志量
		//
		// 默认值: true
		Stat bool `json:",default=true"`

		// KeepDays 日志文件保留天数
		//
		// 仅当Mode为"file"或"volume"时生效
		// 设置为0表示永久保留
		//
		// 示例：
		//  - 7:  保留7天
		//  - 30: 保留30天
		//  - 0:  永久保留（默认）
		KeepDays int `json:",optional"`

		// StackCooldownMillis 堆栈日志记录冷却时间（毫秒）
		//
		// 防止频繁记录堆栈信息影响性能
		// 在指定时间间隔内，相同的堆栈信息只会记录一次
		//
		// 默认值: 100
		StackCooldownMillis int `json:",default=100"`

		// MaxBackups 最大备份日志文件数
		//
		// 仅当Rotation为"size"时生效
		// 设置为0表示无限制
		//
		// 即使设置为0，如果达到KeepDays限制，日志文件仍会被删除
		//
		// 默认值: 0
		MaxBackups int `json:",default=0"`

		// MaxSize 单个日志文件最大大小（MB）
		//
		// 仅当Rotation为"size"时生效
		// 设置为0表示无限制
		//
		// 示例：
		//  - 100: 单个文件最大100MB
		//  - 0:   无限制（默认）
		MaxSize int `json:",default=0"`

		// Rotation 日志轮转规则
		//
		// 可选值：
		//  - "daily": 按天轮转，每天生成新文件
		//  - "size":  按大小轮转，达到MaxSize后生成新文件
		//
		// 默认值: "daily"
		Rotation string `json:",default=daily,options=[daily,size]"`

		// FileTimeFormat 日志文件名中的时间格式
		//
		// 用于按时间分割日志文件时的文件名格式
		// 遵循Go时间格式规范
		//
		// 示例：
		//  - "2006-01-02"        // 按天分割
		//  - "2006-01-02-15"      // 按小时分割
		//
		// 默认值: "2006-01-02T15:04:05.000Z07:00"
		FileTimeFormat string `json:",optional"`

		// FieldKeys 日志字段键名配置
		//
		// 用于自定义日志字段的键名，适配不同的日志收集系统
		// 如ELK、Splunk、Loki等
		FieldKeys fieldKeyConf `json:",optional"`
	}

	// fieldKeyConf 定义日志字段的键名配置
	//
	// 通过自定义键名，可以适配不同的日志分析系统：
	//  - ELK Stack: 通常使用 @timestamp, level, message 等
	//  - Splunk:    可能有不同的字段命名约定
	//  - 自定义系统: 根据需求调整字段名
	fieldKeyConf struct {
		// CallerKey 调用者信息字段键名
		//
		// 包含文件名和行号，用于定位日志输出位置
		//
		// 默认值: "caller"
		CallerKey string `json:",default=caller"`

		// ContentKey 日志内容字段键名
		//
		// 存储主要的日志消息内容
		//
		// 默认值: "content"
		ContentKey string `json:",default=content"`

		// DurationKey 耗时字段键名
		//
		// 用于记录操作耗时，单位通常为毫秒或纳秒
		//
		// 默认值: "duration"
		DurationKey string `json:",default=duration"`

		// LevelKey 日志级别字段键名
		//
		// 存储日志级别信息（debug, info, error等）
		//
		// 默认值: "level"
		LevelKey string `json:",default=level"`

		// SpanKey 分布式追踪跨度字段键名
		//
		// 用于微服务架构中的请求追踪
		//
		// 默认值: "span"
		SpanKey string `json:",default=span"`

		// TimestampKey 时间戳字段键名
		//
		// 存储日志记录时间，默认兼容Elasticsearch格式
		//
		// 默认值: "@timestamp"
		TimestampKey string `json:",default=@timestamp"`

		// TraceKey 分布式追踪ID字段键名
		//
		// 用于标识整个请求链路的唯一ID
		//
		// 默认值: "trace"
		TraceKey string `json:",default=trace"`

		// TruncatedKey 内容截断标记字段键名
		//
		// 当日志内容被截断时，此字段标记为true
		//
		// 默认值: "truncated"
		TruncatedKey string `json:",default=truncated"`
	}
)
