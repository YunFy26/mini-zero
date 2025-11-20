日志模块
一、需求分析：为什么需要日志系统？日志系统要解决的核心问题？
简单来说就一句话：让软件系统的运行“可观测”。
- 问题排查与诊断
  - 快速定位错误：当程序出现异常、警告或错误时，日志能提供最直接的上下文信息，如错误堆栈、当时的变量状态等，让开发者能快速定位到问题根源
  - 分析逻辑漏洞：有些Bug不是错误异常，而是业务逻辑问题。通过记录关键业务流程的日志，可以一步步回溯，分析逻辑是否正确
  - 性能分析：通过记录关键操作（如数据库查询、API调用）的耗时，可以分析出系统的性能瓶颈在哪里
- 行为追踪与审计
  - 安全审计：记录用户的敏感操作，如登陆、修改密码、资金变动、数据访问等。发生安全事件可以通过日志追溯操作链条
  - 用户行为分析：分析用户的操作路径，用于优化产品体验
- 系统监控与运营
  - 监控与告警：通过实时分析日志（例如，统计“ERROR”级别日志在5分钟内出现的频率），可以在系统出现异常时自动触发告警，让运维团队在用户感知之前发现问题
  - 大数据分析：将日志收集到大数据平台（如Elasticsearch、Spark），可以进行深度分析，例如：分析网站流量趋势、发现潜在的安全攻击模式（如某个IP在短时间内大量尝试登录）
  - 数据统计：生成业务报表，如日活用户数、交易量统计等

1. 日志是给谁看的？
  - 开发人员：调试问题、理解程序执行流程
  - 运维人员：监控系统状态、定位生产问题
  - 机器系统：日志采集、分析、告警
  所以，设计决策：
  - 支持人类可读格式 + 机器可读格式
  - 需要结构化日志，方便查询和聚合
2. 日志在什么场景下使用？
  - 开发环境：控制台输出，实时查看
  - 测试环境：文件输出，便于问题复现
  - 生产环境：
    - 高并发
    - 7✖️24小时运行
    - 多实例部署
    - 需要日志轮转、压缩、清理
  所以，设计决策：
  - 必须要高性能（异步写入、缓冲）
  - 必须并发安全
  - 必须零停机（不能阻塞业务）
3. 日志可能遇到什么问题？
  - 日志量太大导致磁盘爆满 -> 日志轮转、限制大小
  - 频繁打印日志影响性能 -> 日志级别控制、异步写入
  - 错误风暴，堆栈刷屏 -> 堆栈日志限流
  - 敏感信息泄漏 -> 日志脱敏
  - 分布式系统，日志分散 -> 链路追踪


我们希望通过日志能够观测程序的运行，日志中的所有内容都需要输出吗？
根据需求输出日志
DebugLevel：输出所有日志
InfoLevel：输出业务日志
ErrorLevel：只输出错误
SevereLevel：只输出严重错误

日志应该包含哪些信息呢？
```go
// 定义日志包含的字段 - 系统字段
const (
    defaultTimestampKey = "@timestamp"  // 时间：何时发生
    defaultLevelKey     = "level"       // 级别：严重程度
    defaultContentKey   = "content"     // 内容：发生了什么
    defaultCallerKey    = "caller"      // 调用者：哪里打印的日志
    defaultTraceKey     = "trace"       // 追踪ID：分布式追踪
    defaultSpanKey      = "span"        // Span ID：调用链
    defaultDurationKey  = "duration"    // 耗时：性能分析
    defaultTruncatedKey = "truncated"   // 截断标记：内容是否完整
)

// 用户也可以自己定义要记录的字段和字段值
type LogField struct {
    Key   string  // 用户指定字段名
    Value any     // 用户指定字段值（任意类型）
}
```
比如开发者要记录下登陆行为，添加了登陆行为需要记录下来的字段名
```go
logx.Infow("user login",
    logx.Field("user_id", 12345),           ← LogField{Key: "user_id", Value: 12345}
    logx.Field("username", "alice"),        ← LogField{Key: "username", Value: "alice"}
    logx.Field("password", "password"),        ← LogField{Key: "password", Value: "password"}
    logx.Field("login_time", time.Now()),   ← LogField{Key: "login_time", Value: time.Now()}
)
```
此时日志包含的信息为：
```json
{
  "@timestamp": "2025-11-19T07:29:15.123+08:00",  ← 系统字段
  "level": "info",                                 ← 系统字段
  "content": "user login",                         ← 系统字段（用户指定了是user login行为）
  "caller": "handler/user.go:45",                 ← 系统字段
  ...
  "user_id": 12345,                                ← 用户字段（LogField）
  "username": "alice",                             ← 用户字段（LogField）
  "password": "password",                          ← 用户字段（LogField）
  "login_time": "2025-11-19T07:29:15.456+08:00"  ← 用户字段（LogField）
}
```
很明显，我们需要把系统定义的字段和用户希望记录的字段统一写入日志，可以创建一个map来达到我们的目的，为这个map定义一个类型：日志条目容器（logEntry）
```go
// 把系统字段和用户字段合并到一起
type logEntry map[string]any

// 在输出日志时，把用户字段和系统字段先放到logEntry里，然后再序列化输出
func output(writer io.Writer, level string, val any, fields ...LogField) {
    ...
    // 1. 创建日志条目容器
    // +3 是系统字段 @timestamp level content
    entry := make(logEntry, len(fields)+3)  
    
    // 2. 把用户的自定义字段放进去
    for _, field := range fields {
        entry[field.Key] = field.Value
    }
    
    // 3. 把系统字段放进去
    entry[timestampKey] = getTimestamp()  // entry["@timestamp"] = "2025-11-19..."
    entry[levelKey] = level               // entry["level"] = "info"
    entry[contentKey] = val               // entry["content"] = "user login"
    ...
}

// 在输出日志时要考虑：是否包含了太长的内容，是否包含了敏感信息

```

日志需要输出到哪里呢？要保留多久呢？单个日志文件最大多大呢？ -》 日志配置
```go
type LogConf struct {
    ServiceName      string  // 服务名（Kubernetes volume模式）
    Mode             string  // console/file/volume
    Encoding         string  // json/plain
    Level            string  // debug/info/error/severe
    Path             string  // 日志文件路径
    
    // 内容控制
    TimeFormat       string  // 时间格式
    MaxContentLength uint32  // 单条日志最大长度
    
    // 轮转控制
    Rotation         string  // daily/size
    KeepDays         int     // 保留天数
    MaxSize          int     // 单文件最大大小（MB）
    MaxBackups       int     // 最大备份数
    Compress         bool    // 是否压缩
    
    // 性能控制
    StackCooldownMillis int  // 堆栈日志冷却时间
    Stat                bool // 是否记录统计日志
    
    // 自定义字段名
    FieldKeys fieldKeyConf
}
```


API设计
1. 对外API：Logger接口
```go
type Logger interface {
    // 五种输出方式（针对每个级别）
    // 1. 简单输出
    // logx.Debug("user:", user)
    Debug(...any)      
    
    // 2. 格式化输出
    // logx.Debugf("user %s login", name)
    Debugf(string, ...any)

    // 3. 延迟计算（性能优化）
    // logx.Debugfn(func() any { return expensiveFunc() })
    Debugfn(func() any)        

    // 4. 对象输出（JSON序列化）
    // logx.Debugv(userStruct)
    Debugv(any)           
    
    // 5. 带字段输出（结构化日志）
    // logx.Debugw("login", Field("user", name))
    Debugw(string, ...LogField) 
    
    // Info/Error/Slow 同样的五种方式
    // ...
    
    // 上下文增强
    WithCallerSkip(skip int) Logger
    WithContext(ctx context.Context) Logger
    WithDuration(d time.Duration) Logger
    WithFields(fields ...LogField) Logger
}
```

API设计
1. 初始化日志系统API
日志输出到哪里？（控制台/文件/Kubernetes）
什么级别的日志需要记录？（Debug/Info/Error）
日志文件如何管理？（轮转/压缩/清理）
将上述配置项封装到LogConf结构体中，通过SetUp或MustSetup方法进行初始化
```go
// 返回error，需要自己处理错误
SetUp(c LogConf) error
// 直接panic，初始化失败直接退出
MustSetup(c LogConf)
```

eg：
```go
func main() {
  logx.MustSetup(logx.LogConf{
    Mode:     "console",
    Encoding: "plain",  // 纯文本，易读
    Level:    "debug",  // 显示所有日志
  })

  // 生产环境：写入文件
  logx.MustSetup(logx.LogConf{
    Mode:     "file",
    Encoding: "json",   // JSON 格式，便于采集
    Level:    "info",   // 只记录重要日志
    Path:     "/var/log/myapp",
  })

  // Kubernetes 环境：区分多个 Pod
  logx.MustSetup(logx.LogConf{
    ServiceName: "user-service",  // 必须指定服务名
    Mode:        "volume",
    Path:        "/var/log",
  })
}
```

有时候运行时需要动态调整日志级别，所以还需要运行时调整API
```go
// ==================== 运行时调整 API ====================

// 动态调整日志级别（不需要重启）
logx.SetLevel(logx.DebugLevel)

// 禁用所有日志（性能测试时）
logx.Disable()

// 禁用统计日志
logx.DisableStat()
```

2. 记录日志API
这条日志是什么级别？（Debug/Info/Error）
怎么打最方便？（简单/格式化/结构化）
需要携带上下文吗？（trace_id、user_id）
创建Logger接口，提供多种日志输出方法
```go
type Logger interface {
    // 基础日志方法
    Debug/Debugf/Debugv/Debugw/Debugfn
    Info/Infof/Infov/Infow/Infofn
    Error/Errorf/Errorv/Errorw/Errorfn
    Slow/Slowf/Slowv/Sloww/Slowfn
    
    // 复杂日志方法
    WithContext(ctx context.Context) Logger
    WithDuration(d time.Duration) Logger
    WithFields(fields ...LogField) Logger
    WithCallerSkip(skip int) Logger
}
```
用户调用logx.Info()等方法时，是怎么被记录下来的呢？


3. 全局增强API
有时候我们希望在日志中携带一些全局信息
```go
// AddGlobalFields 添加全局字段（所有日志都携带）
func AddGlobalFields(fields ...LogField)

// ContextWithFields 在 context 中添加字段
func ContextWithFields(ctx context.Context, fields ...LogField) context.Context
```

4. 日志写入器API

5. 







完整的log调用
```code
用户代码
  │
  │  logx.Info("user login")  ← main.go:15
  │
  └──────────────────────────────────────────────┐
                                                 │
┌────────────────────────────────────────────────▼──────────┐
│ 第1层：包级 API 入口（logs.go）                             │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│ logs.go:196-200                                            │
│ func Info(v ...any) {                                      │
│     if shallLog(InfoLevel) {                              │
│         writeInfo(fmt.Sprint(v...))  ← "user login"       │
│     }                                                      │
│ }                                                          │
└────────────────────────────────┬───────────────────────────┘
                                 │
┌────────────────────────────────▼───────────────────────────┐
│ 第2层：内部写入函数（logs.go）                              │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│ logs.go:575-577                                            │
│ func writeInfo(val any, fields ...LogField) {             │
│     getWriter().Info(val, mergeGlobalFields(addCaller(fields...))...) │
│ }                                                          │
│                                                            │
│ 执行顺序（从内到外）：                                      │
│                                                            │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤2.1：addCaller(fields...)                         │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ logs.go:422-424                                       │ │
│ │ func addCaller(fields ...LogField) []LogField {      │ │
│ │     return append(fields, Field(callerKey, getCaller(callerDepth))) │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用                                                │ │
│ │                                                       │ │
│ │ util.go:10-17                                         │ │
│ │ func getCaller(callDepth int) string {               │ │
│ │     _, file, line, ok := runtime.Caller(callDepth)   │ │
│ │     if !ok {                                          │ │
│ │         return ""                                     │ │
│ │     }                                                 │ │
│ │     return prettyCaller(file, line)                  │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用                                                │ │
│ │                                                       │ │
│ │ util.go:23-35                                         │ │
│ │ func prettyCaller(file string, line int) string {    │ │
│ │     idx := strings.LastIndexByte(file, '/')          │ │
│ │     if idx < 0 {                                      │ │
│ │         return fmt.Sprintf("%s:%d", file, line)      │ │
│ │     }                                                 │ │
│ │     idx = strings.LastIndexByte(file[:idx], '/')     │ │
│ │     if idx < 0 {                                      │ │
│ │         return fmt.Sprintf("%s:%d", file, line)      │ │
│ │     }                                                 │ │
│ │     return fmt.Sprintf("%s:%d", file[idx+1:], line)  │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ 输出：[{Key:"caller", Value:"main.go:15"}]            │ │
│ └──────────────────────────────────────────────────────┘ │
│                          ↓                                 │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤2.2：mergeGlobalFields(...)                       │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ writer.go:353-365                                     │ │
│ │ func mergeGlobalFields(fields []LogField) []LogField { │ │
│ │     globals := globalFields.Load()                    │ │
│ │     if globals == nil {                               │ │
│ │         return fields                                 │ │
│ │     }                                                 │ │
│ │     gf := globals.([]LogField)                        │ │
│ │     ret := make([]LogField, 0, len(gf)+len(fields))  │ │
│ │     ret = append(ret, gf...)        // 全局字段       │ │
│ │     ret = append(ret, fields...)    // 当前字段       │ │
│ │     return ret                                        │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ 输出：[                                                │ │
│ │     {Key:"service", Value:"user-service"},  // 全局   │ │
│ │     {Key:"version", Value:"v1.2.3"},        // 全局   │ │
│ │     {Key:"region", Value:"cn-north-1"},     // 全局   │ │
│ │     {Key:"caller", Value:"main.go:15"},     // 当前   │ │
│ │ ]                                                     │ │
│ └──────────────────────────────────────────────────────┘ │
│                          ↓                                 │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤2.3：getWriter()                                  │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ logs.go:469-476                                       │ │
│ │ func getWriter() Writer {                             │ │
│ │     w := writer.Load()                                │ │
│ │     if w == nil {                                     │ │
│ │         w = writer.StoreIfNil(newConsoleWriter())    │ │
│ │     }                                                 │ │
│ │     return w                                          │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ 输出：Writer 对象（concreteWriter）                    │ │
│ └──────────────────────────────────────────────────────┘ │
│                          ↓                                 │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤2.4：writer.Info(val, fields...)                 │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ writer.go:280-282                                     │ │
│ │ func (w *concreteWriter) Info(v any, fields ...LogField) { │ │
│ │     output(w.infoLog, levelInfo, v, fields...)        │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ 调用参数：                                             │ │
│ │   - w.infoLog: logs/access.log 的 Writer             │ │
│ │   - levelInfo: "info"                                 │ │
│ │   - v: "user login"                                   │ │
│ │   - fields: [{service,...},{version,...},{region,...},{caller,...}] │ │
│ └──────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────┘
                                 │
┌────────────────────────────────▼───────────────────────────┐
│ 第3层：格式化与输出层（writer.go）                          │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│ writer.go:367-399                                          │
│ func output(writer io.Writer, level string, val any, fields ...LogField) { │
│                                                            │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤3.1：内容脱敏和截断                               │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ writer.go:368-378                                     │ │
│ │ switch v := val.(type) {                              │ │
│ │ case string:                                          │ │
│ │     maxLen := atomic.LoadUint32(&maxContentLength)   │ │
│ │     if maxLen > 0 && len(v) > int(maxLen) {          │ │
│ │         val = v[:maxLen]                              │ │
│ │         fields = append(fields, truncatedField)      │ │
│ │     }                                                 │ │
│ │ case Sensitive:                                       │ │
│ │     val = v.MaskSensitive()                          │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ 结果：val = "user login" (无需截断)                   │ │
│ └──────────────────────────────────────────────────────┘ │
│                          ↓                                 │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤3.2：创建日志条目容器                             │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ writer.go:380-381                                     │ │
│ │ // +3 for timestamp, level and content               │ │
│ │ entry := make(logEntry, len(fields)+3)               │ │
│ │                                                       │ │
│ │ 结果：entry = make(map[string]any, 7)                │ │
│ │       // 4个字段 + 3个系统字段                        │ │
│ └──────────────────────────────────────────────────────┘ │
│                          ↓                                 │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤3.3：处理字段值（类型转换+脱敏）                   │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ writer.go:382-387                                     │ │
│ │ for _, field := range fields {                        │ │
│ │     // 3.3.1 先脱敏                                   │ │
│ │     mval := maskSensitive(field.Value)               │ │
│ │                                                       │ │
│ │     // 3.3.2 类型转换                                 │ │
│ │     entry[field.Key] = processFieldValue(mval)       │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用 maskSensitive                                 │ │
│ │                                                       │ │
│ │ sensitive.go:11-17                                    │ │
│ │ func maskSensitive(value any) any {                  │ │
│ │     if s, ok := value.(Sensitive); ok {              │ │
│ │         return s.MaskSensitive()                     │ │
│ │     }                                                 │ │
│ │     return value                                      │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用 processFieldValue                             │ │
│ │                                                       │ │
│ │ writer.go:401-438                                     │ │
│ │ func processFieldValue(value any) any {              │ │
│ │     switch val := value.(type) {                     │ │
│ │     case error:                                       │ │
│ │         return encodeError(val)  // val.Error()      │ │
│ │     case []error:                                     │ │
│ │         var errs []string                             │ │
│ │         for _, err := range val {                     │ │
│ │             errs = append(errs, encodeError(err))    │ │
│ │         }                                             │ │
│ │         return errs                                   │ │
│ │     case time.Duration:                               │ │
│ │         return fmt.Sprint(val)  // "50ms"            │ │
│ │     case []time.Duration:                             │ │
│ │         var durs []string                             │ │
│ │         for _, dur := range val {                     │ │
│ │             durs = append(durs, fmt.Sprint(dur))     │ │
│ │         }                                             │ │
│ │         return durs                                   │ │
│ │     case []time.Time:                                 │ │
│ │         var times []string                            │ │
│ │         for _, t := range val {                       │ │
│ │             times = append(times, fmt.Sprint(t))     │ │
│ │         }                                             │ │
│ │         return times                                  │ │
│ │     case json.Marshaler:                              │ │
│ │         return val  // 保留自定义序列化               │ │
│ │     case fmt.Stringer:                                │ │
│ │         return encodeStringer(val)  // val.String()  │ │
│ │     case []fmt.Stringer:                              │ │
│ │         var strs []string                             │ │
│ │         for _, str := range val {                     │ │
│ │             strs = append(strs, encodeStringer(str)) │ │
│ │         }                                             │ │
│ │         return strs                                   │ │
│ │     default:                                          │ │
│ │         return val  // 原样返回                       │ │
│ │     }                                                 │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ 结果：entry = {                                       │ │
│ │     "service": "user-service",                        │ │
│ │     "version": "v1.2.3",                              │ │
│ │     "region": "cn-north-1",                           │ │
│ │     "caller": "main.go:15",                           │ │
│ │ }                                                     │ │
│ └──────────────────────────────────────────────────────┘ │
│                          ↓                                 │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤3.4：添加系统字段并序列化                         │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ writer.go:389-398                                     │ │
│ │ switch atomic.LoadUint32(&encoding) {                │ │
│ │ case plainEncodingType:                               │ │
│ │     plainFields := buildPlainFields(entry)           │ │
│ │     writePlainAny(writer, level, val, plainFields...) │ │
│ │ default:  // JSON 格式                                │ │
│ │     entry[timestampKey] = getTimestamp()             │ │
│ │     entry[levelKey] = level                          │ │
│ │     entry[contentKey] = val                          │ │
│ │     writeJson(writer, entry)                         │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用 getTimestamp                                  │ │
│ │                                                       │ │
│ │ util.go:19-21                                         │ │
│ │ func getTimestamp() string {                          │ │
│ │     return time.Now().Format(timeFormat)             │ │
│ │ }                                                     │ │
│ │ // timeFormat = "2006-01-02T15:04:05.000Z07:00"      │ │
│ │                                                       │ │
│ │ ↓ 调用 writeJson                                     │ │
│ │                                                       │ │
│ │ writer.go:466-476                                     │ │
│ │ func writeJson(writer io.Writer, info any) {         │ │
│ │     if content, err := marshalJson(info); err != nil { │ │
│ │         log.Printf("err: %s\n\n%s", err.Error(), debug.Stack()) │ │
│ │     } else if writer == nil {                         │ │
│ │         log.Println(string(content))                  │ │
│ │     } else {                                          │ │
│ │         if _, err := writer.Write(append(content, '\n')); err != nil { │ │
│ │             log.Println(err.Error())                  │ │
│ │         }                                             │ │
│ │     }                                                 │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用 marshalJson                                   │ │
│ │                                                       │ │
│ │ writer.go:339-351                                     │ │
│ │ func marshalJson(t interface{}) ([]byte, error) {    │ │
│ │     var buf bytes.Buffer                              │ │
│ │     encoder := json.NewEncoder(&buf)                 │ │
│ │     encoder.SetEscapeHTML(false)  // 不转义 <>&      │ │
│ │     err := encoder.Encode(t)                         │ │
│ │     // 移除 Encoder 自动添加的换行符                  │ │
│ │     if l := buf.Len(); l > 0 && buf.Bytes()[l-1] == '\n' { │ │
│ │         buf.Truncate(l - 1)                           │ │
│ │     }                                                 │ │
│ │     return buf.Bytes(), err                          │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ 结果：JSON字节流                                       │ │
│ │ {"@timestamp":"2025-11-19T08:56:23.123+00:00","level":"info","content":"user login","service":"user-service","version":"v1.2.3","region":"cn-north-1","caller":"main.go:15"} │ │
│ └──────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────┘
                                 │
┌────────────────────────────────▼───────────────────────────┐
│ 第4层：异步写入层（rotatelogger.go）                        │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │
│                                                            │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤4.1：发送到异步通道                               │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ rotatelogger.go:267-275                               │ │
│ │ func (l *RotateLogger) Write(data []byte) (int, error) { │ │
│ │     select {                                          │ │
│ │     case l.channel <- data:  // 非阻塞发送到通道      │ │
│ │         return len(data), nil                         │ │
│ │     case <-l.done:          // 日志系统已关闭         │ │
│ │         log.Println(string(data))                     │ │
│ │         return 0, ErrLogFileClosed                    │ │
│ │     }                                                 │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ 此时用户线程立即返回，继续执行业务逻辑                 │ │
│ │ 总耗时：约 0.7ms                                      │ │
│ └──────────────────────────────────────────────────────┘ │
│                          ↓                                 │
│                    后台 Worker 异步处理                    │
│                          ↓                                 │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤4.2：后台 Worker 接收数据                         │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ rotatelogger.go:376-399                               │ │
│ │ func (l *RotateLogger) startWorker() {                │ │
│ │     l.waitGroup.Add(1)                                │ │
│ │                                                       │ │
│ │     go func() {                                       │ │
│ │         defer l.waitGroup.Done()                     │ │
│ │                                                       │ │
│ │         for {                                         │ │
│ │             select {                                  │ │
│ │             case event := <-l.channel:  // 接收数据   │ │
│ │                 l.write(event)          // 写入文件   │ │
│ │                                                       │ │
│ │             case <-l.done:              // 关闭信号   │ │
│ │                 // 清空通道，确保所有日志都写入       │ │
│ │                 for {                                 │ │
│ │                     select {                          │ │
│ │                     case event := <-l.channel:        │ │
│ │                         l.write(event)                │ │
│ │                     default:                          │ │
│ │                         return  // 通道为空，退出     │ │
│ │                     }                                 │ │
│ │                 }                                     │ │
│ │             }                                         │ │
│ │         }                                             │ │
│ │     }()                                               │ │
│ │ }                                                     │ │
│ └──────────────────────────────────────────────────────┘ │
│                          ↓                                 │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ 步骤4.3：判断轮转 + 写入文件                          │ │
│ │ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ │
│ │ rotatelogger.go:401-414                               │ │
│ │ func (l *RotateLogger) write(v []byte) {              │ │
│ │     // 4.3.1 判断是否需要轮转                         │ │
│ │     if l.rule.ShallRotate(l.currentSize + int64(len(v))) { │ │
│ │         if err := l.rotate(); err != nil {            │ │
│ │             log.Println(err)                          │ │
│ │         } else {                                      │ │
│ │             l.rule.MarkRotated()                     │ │
│ │             l.currentSize = 0                        │ │
│ │         }                                             │ │
│ │     }                                                 │ │
│ │                                                       │ │
│ │     // 4.3.2 写入文件                                 │ │
│ │     if l.fp != nil {                                  │ │
│ │         l.fp.Write(v)                                 │ │
│ │         l.currentSize += int64(len(v))               │ │
│ │     }                                                 │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 如果需要轮转                                        │ │
│ │                                                       │ │
│ │ rotatelogger.go:348-374                               │ │
│ │ func (l *RotateLogger) rotate() error {               │ │
│ │     // 1. 关闭当前文件                                │ │
│ │     if l.fp != nil {                                  │ │
│ │         err := l.fp.Close()                           │ │
│ │         l.fp = nil                                    │ │
│ │         if err != nil {                               │ │
│ │             return err                                │ │
│ │         }                                             │ │
│ │     }                                                 │ │
│ │                                                       │ │
│ │     // 2. 重命名为备份文件                            │ │
│ │     _, err := os.Stat(l.filename)                    │ │
│ │     if err == nil && len(l.backup) > 0 {             │ │
│ │         backupFilename := l.getBackupFilename()      │ │
│ │         err = os.Rename(l.filename, backupFilename)  │ │
│ │         if err != nil {                               │ │
│ │             return err                                │ │
│ │         }                                             │ │
│ │                                                       │ │
│ │         // 3. 异步压缩和清理                          │ │
│ │         l.postRotate(backupFilename)                 │ │
│ │     }                                                 │ │
│ │                                                       │ │
│ │     // 4. 创建新文件                                  │ │
│ │     l.backup = l.rule.BackupFileName()               │ │
│ │     if l.fp, err = os.Create(l.filename); err == nil { │ │
│ │         fs.CloseOnExec(l.fp)                         │ │
│ │     }                                                 │ │
│ │                                                       │ │
│ │     return err                                        │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用 postRotate（异步压缩）                        │ │
│ │                                                       │ │
│ │ rotatelogger.go:340-346                               │ │
│ │ func (l *RotateLogger) postRotate(file string) {     │ │
│ │     go func() {                                       │ │
│ │         l.maybeCompressFile(file)                    │ │
│ │         l.maybeDeleteOutdatedFiles()                 │ │
│ │     }()                                               │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用 maybeCompressFile                             │ │
│ │                                                       │ │
│ │ rotatelogger.go:312-329                               │ │
│ │ func (l *RotateLogger) maybeCompressFile(file string) { │ │
│ │     if !l.compress {                                  │ │
│ │         return                                        │ │
│ │     }                                                 │ │
│ │                                                       │ │
│ │     defer func() {                                    │ │
│ │         if r := recover(); r != nil {                 │ │
│ │             ErrorStack(r)                             │ │
│ │         }                                             │ │
│ │     }()                                               │ │
│ │                                                       │ │
│ │     if _, err := os.Stat(file); err != nil {          │ │
│ │         return                                        │ │
│ │     }                                                 │ │
│ │                                                       │ │
│ │     compressLogFile(file)                            │ │
│ │ }                                                     │ │
│ │                                                       │ │
│ │ ↓ 调用 maybeDeleteOutdatedFiles                      │ │
│ │                                                       │ │
│ │ rotatelogger.go:331-338                               │ │
│ │ func (l *RotateLogger) maybeDeleteOutdatedFiles() {  │ │
│ │     files := l.rule.OutdatedFiles()                  │ │
│ │     for _, file := range files {                      │ │
│ │         if err := os.Remove(file); err != nil {       │ │
│ │             Errorf("failed to remove outdated file: %s", file) │ │
│ │         }                                             │ │
│ │     }                                                 │ │
│ │ }                                                     │ │
│ └──────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
                          logs/access.log
                         (文件系统，磁盘)
```

完整时序图
```code
时间轴                             用户线程                            后台Worker

T0    logx.Info("user login")
      │
T1    ├─ shallLog(InfoLevel) ✓
      │
T2    ├─ fmt.Sprint("user login")
      │
T3    ├─ writeInfo("user login")
      │   ├─ addCaller()
      │   │   └─ getCaller(4)
      │   │       └─ runtime.Caller(4)
      │   │           └─ "main.go:15"
      │   │
      │   ├─ mergeGlobalFields([{caller,...}])
      │   │   └─ 合并全局字段
      │   │       └─ [{service,...},{version,...},{region,...},{caller,...}]
      │   │
      │   ├─ getWriter()
      │   │   └─ 返回 concreteWriter
      │   │
      │   └─ writer.Info(val, fields)
      │
T4    ├─ output(w.infoLog, "info", val, fields)
      │   ├─ 内容脱敏和截断
      │   ├─ 创建 entry (map)
      │   ├─ processFieldValue() 处理每个字段
      │   ├─ 添加系统字段 (timestamp, level, content)
      │   └─ marshalJson(entry)
      │       └─ JSON字节流
      │
T5    └─ rotateLogger.Write(data)
          └─ channel <- data  (非阻塞发送)
              └─ 立即返回 ✓
                                                           │
用户线程继续执行业务逻辑                                     │
                                                           │
                                          T6    ├─ 从通道接收
                                                │  data := <-channel
                                                │
                                          T7    ├─ l.write(data)
                                                │  ├─ ShallRotate?
                                                │  │  └─ 判断文件大小
                                                │  │
                                                │  ├─ 如果需要轮转
                                                │  │  ├─ rotate()
                                                │  │  │  ├─ 关闭当前文件
                                                │  │  │  ├─ 重命名为备份
                                                │  │  │  ├─ 创建新文件
                                                │  │  │  └─ 异步压缩
                                                │  │  │
                                                │  │  └─ MarkRotated()
                                                │  │
                                          T8    │  └─ fp.Write(data)
                                                │     └─ 写入磁盘 ✓
                                                │
                                                ▼
                                          logs/access.log
```