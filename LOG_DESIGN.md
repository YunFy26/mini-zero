# 日志模块完整设计笔记

## 第一章：需求分析 - 为什么需要日志系统？

### 1.1 核心问题
日志系统要解决的核心问题就一句话：**让软件系统的运行"可观测"**。

### 1.2 三大核心需求

#### 1.2.1 问题排查与诊断
- **快速定位错误**：当程序出现异常、警告或错误时，日志能提供最直接的上下文信息，如错误堆栈、当时的变量状态等，让开发者能快速定位到问题根源
- **分析逻辑漏洞**：有些Bug不是错误异常，而是业务逻辑问题。通过记录关键业务流程的日志，可以一步步回溯，分析逻辑是否正确
- **性能分析**：通过记录关键操作（如数据库查询、API调用）的耗时，可以分析出系统的性能瓶颈在哪里

#### 1.2.2 行为追踪与审计
- **安全审计**：记录用户的敏感操作，如登录、修改密码、资金变动、数据访问等。发生安全事件可以通过日志追溯操作链条
- **用户行为分析**：分析用户的操作路径，用于优化产品体验

#### 1.2.3 系统监控与运营
- **监控与告警**：通过实时分析日志（例如，统计"ERROR"级别日志在5分钟内出现的频率），可以在系统出现异常时自动触发告警，让运维团队在用户感知之前发现问题
- **大数据分析**：将日志收集到大数据平台（如Elasticsearch、Spark），可以进行深度分析，例如：分析网站流量趋势、发现潜在的安全攻击模式（如某个IP在短时间内大量尝试登录）
- **数据统计**：生成业务报表，如日活用户数、交易量统计等

### 1.3 关键设计问题

#### 问题1：日志是给谁看的？
- **开发人员**：调试问题、理解程序执行流程
- **运维人员**：监控系统状态、定位生产问题
- **机器系统**：日志采集、分析、告警

**设计决策**：
- 支持人类可读格式（plain text）+ 机器可读格式（JSON）
- 需要结构化日志，方便查询和聚合

#### 问题2：日志在什么场景下使用？
- **开发环境**：控制台输出，实时查看
- **测试环境**：文件输出，便于问题复现
- **生产环境**：
  - 高并发
  - 7×24小时运行
  - 多实例部署
  - 需要日志轮转、压缩、清理

**设计决策**：
- 必须要高性能（异步写入、缓冲）
- 必须并发安全
- 必须零停机（不能阻塞业务）

#### 问题3：日志可能遇到什么问题？
- 日志量太大导致磁盘爆满 → 日志轮转、限制大小
- 频繁打印日志影响性能 → 日志级别控制、异步写入
- 错误风暴，堆栈刷屏 → 堆栈日志限流
- 敏感信息泄漏 → 日志脱敏
- 分布式系统，日志分散 → 链路追踪

---

## 第二章：最小化日志系统设计

### 2.1 设计目标
我们先设计一个**最小化的日志原型系统**，专注于核心功能，暂不考虑性能优化：
1. ✅ 支持不同日志级别（Debug、Info、Error）
2. ✅ 输出到控制台或文件
3. ✅ 记录时间、级别、内容、调用位置
4. ✅ 支持格式化输出
5. ❌ 暂不考虑：异步写入、日志轮转、性能优化

### 2.2 日志级别设计

**为什么需要日志级别？**
不是所有日志都需要输出。例如：
- 开发环境：需要看到所有Debug信息
- 生产环境：只关心Info和Error

```go
// 日志级别定义
const (
    DebugLevel = "debug"  // 详细的调试信息
    InfoLevel  = "info"   // 重要的业务信息
    ErrorLevel = "error"  // 错误信息
)

// 级别优先级：Debug < Info < Error
// 当设置级别为Info时，Debug不会输出，只输出Info和Error
```

### 2.3 日志内容设计

**一条日志应该包含哪些信息？**

#### 2.3.1 系统字段（自动生成）
```go
const (
    timestampKey = "@timestamp"  // 时间：何时发生
    levelKey     = "level"       // 级别：严重程度
    contentKey   = "content"     // 内容：发生了什么
    callerKey    = "caller"      // 调用者：哪里打印的日志
)
```

#### 2.3.2 用户字段（开发者自定义）
```go
// 用户可以添加任意自定义字段
type LogField struct {
    Key   string  // 字段名
    Value any     // 字段值（任意类型）
}
```

#### 2.3.3 示例：用户登录日志
```go
logx.Infow("user login",
    logx.Field("user_id", 12345),
    logx.Field("username", "alice"),
    logx.Field("ip", "192.168.1.100"),
)
```

输出的JSON格式：
```json
{
  "@timestamp": "2025-11-19T12:10:05.123+08:00",
  "level": "info",
  "content": "user login",
  "caller": "handler/user.go:45",
  "user_id": 12345,
  "username": "alice",
  "ip": "192.168.1.100"
}
```

### 2.4 日志条目容器设计

**为什么需要日志条目容器？**
系统字段和用户字段需要合并到一起，然后统一序列化输出。

```go
// logEntry 是一个map，用于存储所有字段
type logEntry map[string]any

// 使用示例
entry := make(logEntry, len(fields)+4) // +4是系统字段数量

// 1. 添加用户字段
for _, field := range fields {
    entry[field.Key] = field.Value
}

// 2. 添加系统字段
entry[timestampKey] = time.Now().Format(time.RFC3339)
entry[levelKey] = "info"
entry[contentKey] = "user login"
entry[callerKey] = "handler/user.go:45"

// 3. 序列化输出
jsonData, _ := json.Marshal(entry)
fmt.Println(string(jsonData))
```

---

## 第三章：系统架构设计

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                   第1层：API接口层                       │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  对外暴露的函数，开发者直接调用                          │
│                                                         │
│  - logx.Debug/Info/Error (v ...any)                   │
│  - logx.Debugf/Infof/Errorf (format string, v ...any) │
│  - logx.Debugw/Infow/Errorw (msg string, fields ...)  │
│  - logx.MustSetup(config LogConf)                      │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ 调用
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   第2层：核心逻辑层                       │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  处理业务逻辑：级别判断、字段合并、调用栈获取             │
│                                                         │
│  - shallLog(level) bool        // 判断是否需要输出       │
│  - addCaller(fields) []Field   // 添加调用位置          │
│  - mergeFields(fields) []Field // 合并全局字段          │
│  - writeInfo/writeError(...)   // 分发到Writer层       │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ 调用
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   第3层：Writer层                        │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  负责日志的格式化和输出控制                              │
│                                                         │
│  - Writer 接口                                          │
│    - Info(v any, fields ...Field)                      │
│    - Error(v any, fields ...Field)                     │
│  - concreteWriter 实现                                  │
│    - infoLog io.Writer   // 普通日志输出               │
│    - errorLog io.Writer  // 错误日志输出               │
│    - output(writer, level, val, fields)  // 核心输出    │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ 调用
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   第4层：格式化层                         │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  处理日志内容：截断、脱敏、类型转换、序列化               │
│                                                         │
│  - 内容截断：限制单条日志最大长度                        │
│  - 内容脱敏：隐藏敏感信息（密码、手机号）                │
│  - 类型转换：error/Duration/Stringer 统一转字符串       │
│  - 序列化：JSON编码或纯文本格式化                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ 写入
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   第5层：输出层                          │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  真正的IO操作：写入控制台或文件                          │
│                                                         │
│  - os.Stdout         // 控制台输出                      │
│  - *os.File          // 文件输出                        │
│  - io.Writer 接口    // 可扩展到任何输出目标             │
└─────────────────────────────────────────────────────────┘
```

### 3.2 为什么要分层？

#### 3.2.1 职责单一
- **API接口层**：只负责接收用户调用，不处理任何业务逻辑
- **核心逻辑层**：只负责判断和决策，不涉及具体的格式化
- **Writer层**：只负责调度输出，不关心日志如何格式化
- **格式化层**：只负责数据转换，不关心输出到哪里
- **输出层**：只负责IO操作，不关心数据内容

#### 3.2.2 易于扩展
- 想增加新的输出格式（如XML）？只需修改格式化层
- 想增加新的输出目标（如网络）？只需修改输出层
- 想增加新的API（如Slow日志）？只需修改API层

#### 3.2.3 易于测试
- 可以单独测试格式化逻辑，无需真实文件
- 可以Mock Writer来测试核心逻辑
- 可以Mock输出层来测试Writer层

---

## 第四章：详细API设计

### 4.1 第1层：API接口层

这一层是**开发者直接接触的部分**，需要简单、直观、易用。

#### 4.1.1 初始化API

**设计思路**：系统启动时必须先配置日志，告诉系统：
- 输出到哪里？（控制台/文件）
- 什么级别？（Debug/Info/Error）
- 什么格式？（JSON/纯文本）

```go
// ==================== 配置结构体 ====================
type LogConf struct {
    Mode     string `json:",default=console,options=[console,file]"`
    Level    string `json:",default=info,options=[debug,info,error]"`
    Encoding string `json:",default=json,options=[json,plain]"`
    Path     string `json:",optional"` // 文件模式时必填
}

// ==================== 初始化API ====================

// MustSetup 初始化日志系统（失败会panic）
// 适用场景：系统启动时，日志必须成功初始化
func MustSetup(c LogConf)

// 使用示例1：开发环境 - 控制台输出
func main() {
    logx.MustSetup(logx.LogConf{
        Mode:     "console",
        Encoding: "plain",  // 纯文本，易读
        Level:    "debug",  // 显示所有日志
    })
}

// 使用示例2：生产环境 - 文件输出
func main() {
    logx.MustSetup(logx.LogConf{
        Mode:     "file",
        Encoding: "json",   // JSON格式，便于采集
        Level:    "info",   // 只记录重要日志
        Path:     "/var/log/myapp/app.log",
    })
}
```

#### 4.1.2 基础日志API

**设计思路**：提供三种输出方式，满足不同场景

```go
// ==================== 方式1：简单输出（拼接多个值）====================
// 适用场景：快速调试
func Debug(v ...any)
func Info(v ...any)
func Error(v ...any)

// 使用示例
user := "alice"
logx.Info("user:", user, "logged in")
// 输出：user: alice logged in

// ==================== 方式2：格式化输出 ====================
// 适用场景：需要控制输出格式
func Debugf(format string, v ...any)
func Infof(format string, v ...any)
func Errorf(format string, v ...any)

// 使用示例
logx.Infof("user %s logged in from %s", "alice", "192.168.1.100")
// 输出：user alice logged in from 192.168.1.100

// ==================== 方式3：结构化输出（带字段）====================
// 适用场景：需要结构化数据，便于后续分析
func Debugw(msg string, fields ...LogField)
func Infow(msg string, fields ...LogField)
func Errorw(msg string, fields ...LogField)

// 使用示例
logx.Infow("user login",
    logx.Field("user_id", 12345),
    logx.Field("username", "alice"),
    logx.Field("ip", "192.168.1.100"),
)
// 输出JSON：
// {
//   "content": "user login",
//   "user_id": 12345,
//   "username": "alice",
//   "ip": "192.168.1.100"
// }
```

#### 4.1.3 字段构造API

```go
// Field 创建一个日志字段
func Field(key string, value any) LogField

// 使用示例
field := logx.Field("user_id", 12345)
// 返回：LogField{Key: "user_id", Value: 12345}
```

### 4.2 第2层：核心逻辑层

这一层是**系统的大脑**，负责决策和调度。

#### 4.2.1 级别判断API

**设计思路**：每次输出日志前，先判断当前级别是否需要输出

```go
// shallLog 判断是否应该输出该级别的日志
// 返回值：true表示应该输出，false表示应该跳过
func shallLog(level string) bool

// 实现逻辑
func shallLog(level string) bool {
    // 获取当前配置的最低级别
    currentLevel := atomic.LoadUint32(&logLevel)
    
    // 将字符串级别转换为数字
    targetLevel := levelToInt(level)
    
    // 只有当目标级别 >= 当前级别时，才输出
    return targetLevel >= currentLevel
}

// 级别转换
func levelToInt(level string) uint32 {
    switch level {
    case DebugLevel:
        return 0
    case InfoLevel:
        return 1
    case ErrorLevel:
        return 2
    default:
        return 0
    }
}

// 使用示例
func Info(v ...any) {
    if shallLog(InfoLevel) {  // 先判断是否需要输出
        writeInfo(fmt.Sprint(v...))
    }
    // 如果当前级别是Error，这里直接返回，不会执行writeInfo
}
```

**为什么要先判断？**
假设生产环境配置级别为Info，但代码中有很多Debug日志：
```go
// 不判断的情况
logx.Debug(expensiveFunction())  // expensiveFunction会执行，浪费性能

// 判断后
logx.Debug(expensiveFunction())
// 内部实现：
if shallLog(DebugLevel) {  // false，直接返回
    // expensiveFunction 不会执行
}
```

#### 4.2.2 调用栈获取API

**设计思路**：自动记录日志是在哪个文件的哪一行打印的

```go
// getCaller 获取调用者的文件名和行号
// callDepth：调用深度，需要跳过多少层
func getCaller(callDepth int) string

// 实现逻辑
func getCaller(callDepth int) string {
    // 使用runtime包获取调用栈信息
    _, file, line, ok := runtime.Caller(callDepth)
    if !ok {
        return ""
    }
    
    // 格式化输出：只保留文件名和行号
    return prettyCaller(file, line)
}

// prettyCaller 格式化调用位置
// /path/to/project/handler/user.go:45 → handler/user.go:45
func prettyCaller(file string, line int) string {
    // 找到最后一个/
    idx := strings.LastIndexByte(file, '/')
    if idx < 0 {
        return fmt.Sprintf("%s:%d", file, line)
    }
    
    // 再往前找一个/，保留两级目录
    idx = strings.LastIndexByte(file[:idx], '/')
    if idx < 0 {
        return fmt.Sprintf("%s:%d", file, line)
    }
    
    return fmt.Sprintf("%s:%d", file[idx+1:], line)
}
```

**调用深度是什么？**
```
用户代码：main.go:15
  │ logx.Info("hello")              ← 这是我们想要的位置
  ▼
API层：logs.go:200
  │ func Info(v ...any) {
  │     writeInfo(fmt.Sprint(v...))
  │ }
  ▼
核心层：logs.go:575
  │ func writeInfo(val any) {
  │     getCaller(4)                 ← callDepth = 4
  │ }
  ▼
runtime.Caller(4)
  │ 跳过4层调用栈
  ▼
返回：main.go:15
```

#### 4.2.3 字段合并API

**设计思路**：将用户字段和系统字段合并

```go
// addCaller 在字段列表中添加调用者信息
func addCaller(fields ...LogField) []LogField

// 实现逻辑
func addCaller(fields ...LogField) []LogField {
    caller := getCaller(callerDepth)
    return append(fields, Field(callerKey, caller))
}

// 使用示例
fields := []LogField{
    Field("user_id", 12345),
    Field("username", "alice"),
}
fields = addCaller(fields...)
// 结果：
// [
//   {Key: "user_id", Value: 12345},
//   {Key: "username", Value: "alice"},
//   {Key: "caller", Value: "main.go:15"}
// ]
```

#### 4.2.4 写入分发API

**设计思路**：根据日志级别，分发到不同的输出目标

```go
// writeInfo 写入Info级别日志
func writeInfo(val any, fields ...LogField)

// writeError 写入Error级别日志
func writeError(val any, fields ...LogField)

// 实现逻辑
func writeInfo(val any, fields ...LogField) {
    // 1. 添加调用者信息
    fields = addCaller(fields...)
    
    // 2. 获取Writer实例
    writer := getWriter()
    
    // 3. 调用Writer的Info方法
    writer.Info(val, fields...)
}

func writeError(val any, fields ...LogField) {
    fields = addCaller(fields...)
    writer := getWriter()
    writer.Error(val, fields...)
}
```

### 4.3 第3层：Writer层

这一层是**调度中心**，负责将日志分发到不同的输出目标。

#### 4.3.1 Writer接口定义

```go
// Writer 日志写入器接口
type Writer interface {
    // Info 写入Info级别日志
    Info(v any, fields ...LogField)
    
    // Error 写入Error级别日志
    Error(v any, fields ...LogField)
    
    // Close 关闭Writer（释放资源）
    Close() error
}
```

**为什么定义接口？**
- 可以有不同的实现（控制台Writer、文件Writer、网络Writer）
- 便于测试（可以Mock一个假的Writer）
- 符合依赖倒置原则

#### 4.3.2 concreteWriter实现

**设计思路**：不同级别的日志可能输出到不同的目标

```go
type concreteWriter struct {
    infoLog  io.Writer  // Info日志输出目标（如 os.Stdout 或文件）
    errorLog io.Writer  // Error日志输出目标（如 os.Stderr 或单独的错误文件）
}

// Info 实现Writer接口的Info方法
func (w *concreteWriter) Info(v any, fields ...LogField) {
    output(w.infoLog, levelInfo, v, fields...)
}

// Error 实现Writer接口的Error方法
func (w *concreteWriter) Error(v any, fields ...LogField) {
    output(w.errorLog, levelError, v, fields...)
}

// Close 关闭Writer（如果是文件，需要关闭文件句柄）
func (w *concreteWriter) Close() error {
    // 最小化版本暂不实现
    return nil
}
```

**为什么Info和Error分开？**
```
控制台模式：
  Info  → os.Stdout（标准输出，绿色）
  Error → os.Stderr（标准错误，红色）

文件模式：
  Info  → logs/access.log（普通日志）
  Error → logs/error.log（错误日志，便于告警）
```

#### 4.3.3 Writer创建API

```go
// getWriter 获取全局Writer实例（单例模式）
func getWriter() Writer

// 实现逻辑
var writer atomic.Value  // 存储Writer实例

func getWriter() Writer {
    w := writer.Load()
    if w == nil {
        // 首次调用，创建默认Writer
        w = writer.StoreIfNil(newConsoleWriter())
    }
    return w.(Writer)
}

// newConsoleWriter 创建控制台Writer
func newConsoleWriter() Writer {
    return &concreteWriter{
        infoLog:  os.Stdout,
        errorLog: os.Stderr,
    }
}

// newFileWriter 创建文件Writer
func newFileWriter(infoPath, errorPath string) (Writer, error) {
    infoFile, err := os.OpenFile(infoPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    errorFile, err := os.OpenFile(errorPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        infoFile.Close()
        return nil, err
    }
    
    return &concreteWriter{
        infoLog:  infoFile,
        errorLog: errorFile,
    }, nil
}
```

### 4.4 第4层：格式化层

这一层是**数据处理中心**，负责将各种类型的数据转换为可输出的格式。

#### 4.4.1 核心输出API

```go
// output 核心输出函数
// writer: 输出目标（os.Stdout或文件）
// level: 日志级别
// val: 日志内容
// fields: 用户自定义字段
func output(writer io.Writer, level string, val any, fields ...LogField)

// 实现逻辑
func output(writer io.Writer, level string, val any, fields ...LogField) {
    // ===== 步骤1：内容截断（防止单条日志过大）=====
    switch v := val.(type) {
    case string:
        maxLen := atomic.LoadUint32(&maxContentLength)
        if maxLen > 0 && len(v) > int(maxLen) {
            val = v[:maxLen]  // 截断
            fields = append(fields, Field(truncatedKey, true))
        }
    }
    
    // ===== 步骤2：创建日志条目容器 =====
    // +3 是系统字段：timestamp、level、content
    entry := make(logEntry, len(fields)+3)
    
    // ===== 步骤3：处理用户字段 =====
    for _, field := range fields {
        entry[field.Key] = processFieldValue(field.Value)
    }
    
    // ===== 步骤4：添加系统字段 =====
    entry[timestampKey] = getTimestamp()
    entry[levelKey] = level
    entry[contentKey] = val
    
    // ===== 步骤5：序列化并输出 =====
    encoding := atomic.LoadUint32(&encoding)
    switch encoding {
    case jsonEncodingType:
        writeJson(writer, entry)
    case plainEncodingType:
        writePlain(writer, entry)
    }
}
```

#### 4.4.2 类型转换API

**设计思路**：Go语言中有很多特殊类型需要特殊处理

```go
// processFieldValue 处理字段值（类型转换）
func processFieldValue(value any) any

// 实现逻辑
func processFieldValue(value any) any {
    switch val := value.(type) {
    // 1. error类型 → 转为字符串
    case error:
        return val.Error()
    
    // 2. time.Duration → 转为可读字符串
    case time.Duration:
        return val.String()  // "50ms", "2s"
    
    // 3. fmt.Stringer接口 → 调用String()方法
    case fmt.Stringer:
        return val.String()
    
    // 4. 其他类型 → 原样返回
    default:
        return val
    }
}

// 使用示例
err := errors.New("connection failed")
processFieldValue(err)  // 返回："connection failed"

duration := 50 * time.Millisecond
processFieldValue(duration)  // 返回："50ms"
```

#### 4.4.3 时间格式化API

```go
// getTimestamp 获取格式化的时间戳
func getTimestamp() string

// 实现逻辑
var timeFormat string = time.RFC3339  // "2006-01-02T15:04:05Z07:00"

func getTimestamp() string {
    return time.Now().Format(timeFormat)
}

// 输出示例
"2025-11-19T12:10:05+08:00"
```

#### 4.4.4 JSON序列化API

```go
// writeJson 将日志条目序列化为JSON并写入
func writeJson(writer io.Writer, entry logEntry)

// 实现逻辑
func writeJson(writer io.Writer, entry logEntry) {
    // 1. 序列化为JSON
    data, err := json.Marshal(entry)
    if err != nil {
        log.Printf("failed to marshal log: %v\n", err)
        return
    }
    
    // 2. 写入目标（加换行符）
    if _, err := writer.Write(append(data, '\n')); err != nil {
        log.Printf("failed to write log: %v\n", err)
    }
}

// 输出示例
{"@timestamp":"2025-11-19T12:10:05+08:00","level":"info","content":"user login","user_id":12345}
```

#### 4.4.5 纯文本格式化API

```go
// writePlain 将日志条目格式化为纯文本并写入
func writePlain(writer io.Writer, entry logEntry)

// 实现逻辑
func writePlain(writer io.Writer, entry logEntry) {
    var buf bytes.Buffer
    
    // 1. 写入时间和级别
    buf.WriteString(entry[timestampKey].(string))
    buf.WriteString(" ")
    buf.WriteString(entry[levelKey].(string))
    buf.WriteString(" ")
    
    // 2. 写入内容
    buf.WriteString(fmt.Sprint(entry[contentKey]))
    
    // 3. 写入其他字段
    for key, value := range entry {
        if key == timestampKey || key == levelKey || key == contentKey {
            continue
        }
        buf.WriteString(" ")
        buf.WriteString(key)
        buf.WriteString("=")
        buf.WriteString(fmt.Sprint(value))
    }
    
    // 4. 写入换行符
    buf.WriteString("\n")
    
    // 5. 输出
    writer.Write(buf.Bytes())
}

// 输出示例
2025-11-19T12:10:05+08:00 info user login user_id=12345 username=alice
```

### 4.5 第5层：输出层

这一层是**IO操作层**，真正执行写入操作。

#### 4.5.1 io.Writer接口

Go标准库的io.Writer接口：
```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

所有输出目标都实现了这个接口：
- `os.Stdout`：控制台标准输出
- `os.Stderr`：控制台标准错误
- `*os.File`：文件
- `bytes.Buffer`：内存缓冲区（测试用）
- 网络连接、HTTP响应等

#### 4.5.2 输出目标

```go
// 控制台输出
os.Stdout.Write([]byte("hello\n"))

// 文件输出
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
file.Write([]byte("hello\n"))
```

---

## 第五章：完整调用流程

### 5.1 流程图

```
用户代码（main.go:15）
  │
  │ logx.Info("user login")
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│ 第1层：API接口层                                         │
├─────────────────────────────────────────────────────────┤
│ logs.go:200                                              │
│ func Info(v ...any) {                                    │
│     if shallLog(InfoLevel) {        ← 判断级别           │
│         writeInfo(fmt.Sprint(v...)) ← "user login"      │
│     }                                                    │
│ }                                                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ 第2层：核心逻辑层                                         │
├─────────────────────────────────────────────────────────┤
│ logs.go:575                                              │
│ func writeInfo(val any, fields ...LogField) {           │
│     fields = addCaller(fields...)  ← 添加调用位置        │
│     // fields = [{caller: "main.go:15"}]                │
│                                                          │
│     writer := getWriter()          ← 获取Writer实例      │
│     writer.Info(val, fields...)    ← 调用Writer         │
│ }                                                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ 第3层：Writer层                                          │
├─────────────────────────────────────────────────────────┤
│ writer.go:280                                            │
│ func (w *concreteWriter) Info(v any, fields ...LogField) { │
│     output(w.infoLog, levelInfo, v, fields...)          │
│     // w.infoLog = os.Stdout 或文件                      │
│ }                                                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ 第4层：格式化层                                           │
├─────────────────────────────────────────────────────────┤
│ writer.go:389                                            │
│ func output(writer io.Writer, level, val, fields) {     │
│     // 1. 创建日志条目                                   │
│     entry := make(logEntry, len(fields)+3)              │
│                                                          │
│     // 2. 处理字段                                       │
│     for _, field := range fields {                       │
│         entry[field.Key] = processFieldValue(field.Value) │
│     }                                                    │
│                                                          │
│     // 3. 添加系统字段                                   │
│     entry["@timestamp"] = getTimestamp()                │
│     entry["level"] = "info"                             │
│     entry["content"] = "user login"                     │
│     entry["caller"] = "main.go:15"                      │
│                                                          │
│     // 4. 序列化                                         │
│     writeJson(writer, entry)                            │
│ }                                                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ 第5层：输出层                                             │
├─────────────────────────────────────────────────────────┤
│ writer.go:466                                            │
│ func writeJson(writer io.Writer, entry logEntry) {      │
│     data, _ := json.Marshal(entry)                      │
│     writer.Write(append(data, '\n'))                    │
│ }                                                        │
│                                                          │
│ 输出到：                                                  │
│ os.Stdout 或 文件                                        │
└─────────────────────────────────────────────────────────┘
```

### 5.2 代码示例

```go
package main

import "logx"

func main() {
    // ===== 步骤1：初始化日志系统 =====
    logx.MustSetup(logx.LogConf{
        Mode:     "console",
        Encoding: "json",
        Level:    "info",
    })
    
    // ===== 步骤2：记录日志 =====
    logx.Info("user login")
    
    // ===== 内部执行流程 =====
    // 1. API层：Info("user login")
    //    └─ shallLog(InfoLevel) → true
    //    └─ writeInfo("user login")
    //
    // 2. 核心层：writeInfo("user login")
    //    └─ addCaller() → [{caller: "main.go:15"}]
    //    └─ getWriter() → concreteWriter实例
    //    └─ writer.Info("user login", [{caller: "main.go:15"}])
    //
    // 3. Writer层：concreteWriter.Info(...)
    //    └─ output(os.Stdout, "info", "user login", [{caller: "main.go:15"}])
    //
    // 4. 格式化层：output(...)
    //    └─ entry = {
    //         "@timestamp": "2025-11-19T12:10:05+08:00",
    //         "level": "info",
    //         "content": "user login",
    //         "caller": "main.go:15"
    //       }
    //    └─ writeJson(os.Stdout, entry)
    //
    // 5. 输出层：writeJson(...)
    //    └─ json.Marshal(entry) → JSON字节流
    //    └─ os.Stdout.Write(JSON字节流)
    //
    // ===== 最终输出到控制台 =====
    // {"@timestamp":"2025-11-19T12:10:05+08:00","level":"info","content":"user login","caller":"main.go:15"}
}
```

### 5.3 时序图

```
时间   用户线程                  系统调用栈
─────────────────────────────────────────────────────────
T0     main.go:15
       logx.Info("user login")
         │
T1       └─→ logs.go:200 Info(v)
               │ shallLog(InfoLevel) → true
               │
T2             └─→ logs.go:575 writeInfo(val)
                     │ addCaller()
                     │   └─→ runtime.Caller(4)
                     │       └─→ "main.go:15"
                     │
T3               │ getWriter()
                 │   └─→ 返回 concreteWriter
                     │
T4                   └─→ writer.go:280 Info(val, fields)
                           │
T5                         └─→ writer.go:389 output(...)
                                 │ 创建 entry
                                 │ 处理字段
                                 │ 添加系统字段
                                 │
T6                               └─→ writer.go:466 writeJson(...)
                                       │ json.Marshal(entry)
                                       │
T7                                     └─→ os.Stdout.Write(data)
                                             │
                                             ▼
                                         控制台输出
```

---

## 第六章：最小化实现代码

### 6.1 目录结构

```
logx/
├── logx.go              // API接口层 + 核心逻辑层
├── writer.go            // Writer层 + 格式化层
├── config.go            // 配置相关
└── types.go             // 类型定义
```

### 6.2 types.go - 类型定义

```go
package logx

// LogField 日志字段
type LogField struct {
    Key   string
    Value any
}

// Field 创建一个日志字段
func Field(key string, value any) LogField {
    return LogField{Key: key, Value: value}
}

// logEntry 日志条目容器
type logEntry map[string]any

// 日志级别常量
const (
    DebugLevel = "debug"
    InfoLevel  = "info"
    ErrorLevel = "error"
)

// 字段名常量
const (
    timestampKey = "@timestamp"
    levelKey     = "level"
    contentKey   = "content"
    callerKey    = "caller"
)

// 编码类型
const (
    jsonEncodingType  = 0
    plainEncodingType = 1
)
```

### 6.3 config.go - 配置

```go
package logx

// LogConf 日志配置
type LogConf struct {
    Mode     string `json:",default=console,options=[console,file]"`
    Level    string `json:",default=info,options=[debug,info,error]"`
    Encoding string `json:",default=json,options=[json,plain]"`
    Path     string `json:",optional"` // 文件模式时必填
}
```

### 6.4 logx.go - API接口层 + 核心逻辑层

```go
package logx

import (
    "fmt"
    "runtime"
    "strings"
    "sync/atomic"
)

// ==================== 全局变量 ====================
var (
    writer    atomic.Value  // Writer实例
    logLevel  atomic.Uint32 // 当前日志级别
    encoding  atomic.Uint32 // 编码类型
    callerDepth = 4         // 调用深度
)

// ==================== 初始化API ====================

// MustSetup 初始化日志系统（失败会panic）
func MustSetup(c LogConf) {
    if err := setupWithConfig(c); err != nil {
        panic(err)
    }
}

func setupWithConfig(c LogConf) error {
    // 1. 设置日志级别
    setLevel(c.Level)
    
    // 2. 设置编码类型
    setEncoding(c.Encoding)
    
    // 3. 创建Writer
    var w Writer
    var err error
    
    switch c.Mode {
    case "console":
        w = newConsoleWriter()
    case "file":
        if c.Path == "" {
            return fmt.Errorf("file mode requires path")
        }
        w, err = newFileWriter(c.Path)
        if err != nil {
            return err
        }
    default:
        return fmt.Errorf("unknown mode: %s", c.Mode)
    }
    
    writer.Store(w)
    return nil
}

func setLevel(level string) {
    switch level {
    case DebugLevel:
        logLevel.Store(0)
    case InfoLevel:
        logLevel.Store(1)
    case ErrorLevel:
        logLevel.Store(2)
    default:
        logLevel.Store(1) // 默认Info
    }
}

func setEncoding(enc string) {
    switch enc {
    case "json":
        encoding.Store(jsonEncodingType)
    case "plain":
        encoding.Store(plainEncodingType)
    default:
        encoding.Store(jsonEncodingType) // 默认JSON
    }
}

// ==================== 基础日志API ====================

// Debug 输出Debug级别日志
func Debug(v ...any) {
    if shallLog(DebugLevel) {
        writeDebug(fmt.Sprint(v...))
    }
}

// Debugf 格式化输出Debug级别日志
func Debugf(format string, v ...any) {
    if shallLog(DebugLevel) {
        writeDebug(fmt.Sprintf(format, v...))
    }
}

// Debugw 带字段输出Debug级别日志
func Debugw(msg string, fields ...LogField) {
    if shallLog(DebugLevel) {
        writeDebug(msg, fields...)
    }
}

// Info 输出Info级别日志
func Info(v ...any) {
    if shallLog(InfoLevel) {
        writeInfo(fmt.Sprint(v...))
    }
}

// Infof 格式化输出Info级别日志
func Infof(format string, v ...any) {
    if shallLog(InfoLevel) {
        writeInfo(fmt.Sprintf(format, v...))
    }
}

// Infow 带字段输出Info级别日志
func Infow(msg string, fields ...LogField) {
    if shallLog(InfoLevel) {
        writeInfo(msg, fields...)
    }
}

// Error 输出Error级别日志
func Error(v ...any) {
    if shallLog(ErrorLevel) {
        writeError(fmt.Sprint(v...))
    }
}

// Errorf 格式化输出Error级别日志
func Errorf(format string, v ...any) {
    if shallLog(ErrorLevel) {
        writeError(fmt.Sprintf(format, v...))
    }
}

// Errorw 带字段输出Error级别日志
func Errorw(msg string, fields ...LogField) {
    if shallLog(ErrorLevel) {
        writeError(msg, fields...)
    }
}

// ==================== 核心逻辑层 ====================

// shallLog 判断是否应该输出该级别的日志
func shallLog(level string) bool {
    currentLevel := logLevel.Load()
    targetLevel := levelToUint32(level)
    return targetLevel >= currentLevel
}

func levelToUint32(level string) uint32 {
    switch level {
    case DebugLevel:
        return 0
    case InfoLevel:
        return 1
    case ErrorLevel:
        return 2
    default:
        return 0
    }
}

// writeDebug 写入Debug级别日志
func writeDebug(val any, fields ...LogField) {
    fields = addCaller(fields...)
    getWriter().Debug(val, fields...)
}

// writeInfo 写入Info级别日志
func writeInfo(val any, fields ...LogField) {
    fields = addCaller(fields...)
    getWriter().Info(val, fields...)
}

// writeError 写入Error级别日志
func writeError(val any, fields ...LogField) {
    fields = addCaller(fields...)
    getWriter().Error(val, fields...)
}

// addCaller 添加调用者信息
func addCaller(fields ...LogField) []LogField {
    return append(fields, Field(callerKey, getCaller(callerDepth)))
}

// getCaller 获取调用者的文件名和行号
func getCaller(callDepth int) string {
    _, file, line, ok := runtime.Caller(callDepth)
    if !ok {
        return ""
    }
    return prettyCaller(file, line)
}

// prettyCaller 格式化调用位置
func prettyCaller(file string, line int) string {
    idx := strings.LastIndexByte(file, '/')
    if idx < 0 {
        return fmt.Sprintf("%s:%d", file, line)
    }
    idx = strings.LastIndexByte(file[:idx], '/')
    if idx < 0 {
        return fmt.Sprintf("%s:%d", file, line)
    }
    return fmt.Sprintf("%s:%d", file[idx+1:], line)
}

// getWriter 获取Writer实例
func getWriter() Writer {
    w := writer.Load()
    if w == nil {
        // 默认使用控制台输出
        w = newConsoleWriter()
        writer.Store(w)
    }
    return w.(Writer)
}
```

### 6.5 writer.go - Writer层 + 格式化层

```go
package logx

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "os"
    "time"
)

// ==================== Writer接口 ====================

// Writer 日志写入器接口
type Writer interface {
    Debug(v any, fields ...LogField)
    Info(v any, fields ...LogField)
    Error(v any, fields ...LogField)
    Close() error
}

// ==================== concreteWriter实现 ====================

type concreteWriter struct {
    infoLog  io.Writer
    errorLog io.Writer
}

func (w *concreteWriter) Debug(v any, fields ...LogField) {
    output(w.infoLog, DebugLevel, v, fields...)
}

func (w *concreteWriter) Info(v any, fields ...LogField) {
    output(w.infoLog, InfoLevel, v, fields...)
}

func (w *concreteWriter) Error(v any, fields ...LogField) {
    output(w.errorLog, ErrorLevel, v, fields...)
}

func (w *concreteWriter) Close() error {
    return nil
}

// ==================== Writer创建 ====================

// newConsoleWriter 创建控制台Writer
func newConsoleWriter() Writer {
    return &concreteWriter{
        infoLog:  os.Stdout,
        errorLog: os.Stderr,
    }
}

// newFileWriter 创建文件Writer
func newFileWriter(path string) (Writer, error) {
    file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    return &concreteWriter{
        infoLog:  file,
        errorLog: file,
    }, nil
}

// ==================== 格式化层 ====================

// output 核心输出函数
func output(writer io.Writer, level string, val any, fields ...LogField) {
    // 1. 创建日志条目容器
    entry := make(logEntry, len(fields)+3)
    
    // 2. 处理用户字段
    for _, field := range fields {
        entry[field.Key] = processFieldValue(field.Value)
    }
    
    // 3. 添加系统字段
    entry[timestampKey] = getTimestamp()
    entry[levelKey] = level
    entry[contentKey] = val
    
    // 4. 序列化输出
    enc := encoding.Load()
    switch enc {
    case jsonEncodingType:
        writeJson(writer, entry)
    case plainEncodingType:
        writePlain(writer, entry)
    }
}

// processFieldValue 处理字段值（类型转换）
func processFieldValue(value any) any {
    switch val := value.(type) {
    case error:
        return val.Error()
    case time.Duration:
        return val.String()
    case fmt.Stringer:
        return val.String()
    default:
        return val
    }
}

// getTimestamp 获取格式化的时间戳
func getTimestamp() string {
    return time.Now().Format(time.RFC3339)
}

// writeJson JSON序列化输出
func writeJson(writer io.Writer, entry logEntry) {
    data, err := json.Marshal(entry)
    if err != nil {
        log.Printf("failed to marshal log: %v\n", err)
        return
    }
    
    if _, err := writer.Write(append(data, '\n')); err != nil {
        log.Printf("failed to write log: %v\n", err)
    }
}

// writePlain 纯文本格式化输出
func writePlain(writer io.Writer, entry logEntry) {
    var buf bytes.Buffer
    
    // 时间 级别 内容
    buf.WriteString(entry[timestampKey].(string))
    buf.WriteString(" ")
    buf.WriteString(entry[levelKey].(string))
    buf.WriteString(" ")
    buf.WriteString(fmt.Sprint(entry[contentKey]))
    
    // 其他字段
    for key, value := range entry {
        if key == timestampKey || key == levelKey || key == contentKey {
            continue
        }
        buf.WriteString(" ")
        buf.WriteString(key)
        buf.WriteString("=")
        buf.WriteString(fmt.Sprint(value))
    }
    
    buf.WriteString("\n")
    writer.Write(buf.Bytes())
}
```

---

## 第七章：使用示例

### 7.1 控制台输出（JSON格式）

```go
package main

import "logx"

func main() {
    // 初始化
    logx.MustSetup(logx.LogConf{
        Mode:     "console",
        Encoding: "json",
        Level:    "debug",
    })
    
    // 简单输出
    logx.Info("server started")
    
    // 格式化输出
    logx.Infof("listening on port %d", 8080)
    
    // 结构化输出
    logx.Infow("user login",
        logx.Field("user_id", 12345),
        logx.Field("username", "alice"),
    )
}
```

**输出：**
```json
{"@timestamp":"2025-11-19T12:10:05+08:00","level":"info","content":"server started","caller":"main.go:10"}
{"@timestamp":"2025-11-19T12:10:05+08:00","level":"info","content":"listening on port 8080","caller":"main.go:13"}
{"@timestamp":"2025-11-19T12:10:05+08:00","level":"info","content":"user login","caller":"main.go:16","user_id":12345,"username":"alice"}
```

### 7.2 控制台输出（纯文本格式）

```go
logx.MustSetup(logx.LogConf{
    Mode:     "console",
    Encoding: "plain",
    Level:    "info",
})

logx.Infow("user login",
    logx.Field("user_id", 12345),
    logx.Field("username", "alice"),
)
```

**输出：**
```
2025-11-19T12:10:05+08:00 info user login caller=main.go:10 user_id=12345 username=alice
```

### 7.3 文件输出

```go
logx.MustSetup(logx.LogConf{
    Mode:     "file",
    Encoding: "json",
    Level:    "info",
    Path:     "/var/log/myapp/app.log",
})

logx.Info("server started")
```

**app.log内容：**
```json
{"@timestamp":"2025-11-19T12:10:05+08:00","level":"info","content":"server started","caller":"main.go:10"}
```

### 7.4 级别过滤

```go
logx.MustSetup(logx.LogConf{
    Level: "info",  // 只输出Info及以上级别
})

logx.Debug("this is debug")  // 不会输出
logx.Info("this is info")     // 会输出
logx.Error("this is error")   // 会输出
```

---