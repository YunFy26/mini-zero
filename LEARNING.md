# mini-zero 学习计划与日志

## 🎯 学习目标

通过逐步实现 go-zero 的核心功能，深入理解微服务框架设计和 Go 语言最佳实践。本计划基于 [go-zero](https://github.com/zeromicro/go-zero) 官方仓库的完整架构制定。

## 📅 全面学习计划

> 参考：go-zero 是一个集成了各种工程实践的 web 和 rpc 框架，包含简洁的 API 定义语法和强大的代码生成工具 goctl。

---

### 🏗️ 第一阶段：Core 核心包（进行中）

#### 1.1 日志系统 (core/logx) ⭐ 优先级：高
- [x] Logger 接口设计
- [x] 配置管理 (LogConf)
- [x] 延迟求值优化 (Debugfn)
- [x] 颜色工具 (color)
- [ ] 日志级别完整实现 (Debug/Info/Error/Severe/Slow/Stack/Stat)
- [ ] 多种 Writer 实现 (console/file/rotate)
- [ ] 日志轮转 (rotation)
- [ ] 结构化日志 (JSON/Plain encoding)
- [ ] 异步写入优化
- [ ] 全局字段 (Global Fields)
- [ ] Context 集成
- [ ] 日志采样 (Sampling)

#### 1.2 并发控制 (core/syncx) ⭐ 优先级：高
- [x] AtomicBool
- [ ] AtomicDuration
- [ ] AtomicFloat64
- [ ] Barrier (屏障)
- [ ] Cond (条件变量)
- [ ] DoneChan
- [ ] ImmutableResource (不可变资源)
- [ ] Limit (限制器)
- [ ] ManagedResource (托管资源)
- [ ] Once (单次执行)
- [ ] Pool (对象池)
- [ ] ResourceManager (资源管理器)
- [ ] SharedCalls (合并调用)
- [ ] TimeoutLimit (超时限制)

#### 1.3 并发模式 (core/mr) ⭐ 优先级：中
- [ ] MapReduce 实现
- [ ] FinishVoid (完成无返回)
- [ ] Finish (完成有返回)
- [ ] Map (映射)
- [ ] MapVoid (映射无返回)

#### 1.4 集合工具 (core/collection) ⭐ 优先级：中
- [ ] Cache (缓存)
- [ ] LRU Cache
- [ ] Ring (环形缓冲)
- [ ] RollingWindow (滚动窗口)
- [ ] SafeMap (安全 Map)
- [ ] Set (集合)
- [ ] TimingWheel (时间轮)

#### 1.5 限流器 (core/limit) ⭐ 优先级：高
- [ ] PeriodLimit (周期限流)
- [ ] TokenBucket (令牌桶)

#### 1.6 熔断器 (core/breaker) ⭐ 优先级：高
- [ ] CircuitBreaker 接口
- [ ] Google SRE Breaker
- [ ] 熔断器状态管理
- [ ] 自适应熔断

#### 1.7 自适应降载 (core/load) ⭐ 优先级：中
- [ ] AdaptiveShedder (自适应降载器)
- [ ] 负载监控
- [ ] 自动降载策略

#### 1.8 统计监控 (core/stat) ⭐ 优先级：中
- [ ] Metrics 接口
- [ ] Usage (资源使用统计)
- [ ] RemoteWriter (远程写入)

#### 1.9 链路追踪 (core/trace) ⭐ 优先级：中
- [ ] Tracer 接口
- [ ] OpenTelemetry 集成
- [ ] Span 管理
- [ ] TraceID/SpanID 生成

#### 1.10 配置管理 (core/conf) ⭐ 优先级：中
- [ ] 配置加载
- [ ] 环境变量支持
- [ ] 配置热更新

#### 1.11 服务发现 (core/discov) ⭐ 优先级：高
- [ ] Registry 接口
- [ ] etcd 集成
- [ ] Consul 支持
- [ ] Kubernetes 支持

#### 1.12 存储抽象 (core/stores) ⭐ 优先级：高
**1.12.1 Cache 缓存**
- [ ] Redis Cache
- [ ] 本地缓存
- [ ] 缓存策略 (旁路/穿透/击穿/雪崩)

**1.12.2 SQL 数据库**
- [ ] MySQL 集成
- [ ] PostgreSQL 集成
- [ ] 连接池管理
- [ ] 事务支持
- [ ] SQL Builder

**1.12.3 MongoDB**
- [ ] MongoDB 客户端封装
- [ ] CRUD 操作

**1.12.4 Redis**
- [ ] Redis 客户端封装
- [ ] 分布式锁
- [ ] Redis Script

#### 1.13 线程安全 (core/threading) ⭐ 优先级：低
- [ ] GoSafe (安全 goroutine)
- [ ] RunSafe (安全执行)
- [ ] TaskRunner (任务执行器)
- [ ] WorkerPool (工作池)

#### 1.14 时间工具 (core/timex) ⭐ 优先级：低
- [ ] Ticker
- [ ] Time 格式化

#### 1.15 工具函数 (core/stringx, core/mathx, core/fx) ⭐ 优先级：低
- [ ] 字符串工具
- [ ] 数学工具
- [ ] 函数式工具

---

### 🌐 第二阶段：REST 框架 (rest)

#### 2.1 HTTP 服务器核心 ⭐ 优先级：高
- [ ] Server 结构
- [ ] 路由设计 (Router)
- [ ] 路由组 (Group)
- [ ] 参数绑定
- [ ] 请求验证

#### 2.2 中间件系统 ⭐ 优先级：高
- [ ] 中间件接口
- [ ] AuthMiddleware (认证)
- [ ] CorsMiddleware (跨域)
- [ ] LogMiddleware (日志)
- [ ] MetricsMiddleware (指标)
- [ ] PrometheusMiddleware
- [ ] RecoverMiddleware (恢复)
- [ ] TimeoutMiddleware (超时)
- [ ] TraceMiddleware (追踪)
- [ ] BreakerMiddleware (熔断)
- [ ] SheddingMiddleware (降载)
- [ ] MaxConnsMiddleware (最大连接)
- [ ] MaxBytesMiddleware (最大字节)

#### 2.3 响应处理 ⭐ 优先级：中
- [ ] httpx 工具
- [ ] 错误处理
- [ ] 响应格式化

#### 2.4 路由引擎 ⭐ 优先级：中
- [ ] 路由树 (Trie)
- [ ] 路径参数
- [ ] 查询参数

---

### 🔌 第三阶段：RPC 框架 (zrpc)

#### 3.1 gRPC 集成 ⭐ 优先级：高
- [ ] Server 实现
- [ ] Client 实现
- [ ] 拦截器 (Interceptor)

#### 3.2 服务治理 ⭐ 优先级：高
- [ ] 服务注册
- [ ] 服务发现
- [ ] 负载均衡 (P2C/Random/RoundRobin/Weighted)
- [ ] 超时控制
- [ ] 重试机制
- [ ] 熔断降级

#### 3.3 RPC 中间件 ⭐ 优先级：中
- [ ] Auth 拦截器
- [ ] Breaker 拦截器
- [ ] Prometheus 拦截器
- [ ] Timeout 拦截器
- [ ] Trace 拦截器

---

### 🚪 第四阶段：网关 (gateway)

#### 4.1 API 网关 ⭐ 优先级：低
- [ ] HTTP 到 gRPC 转换
- [ ] 路由转发
- [ ] 协议转换

---

### 🛠️ 第五阶段：工具链 (tools/goctl)

#### 5.1 代码生成 ⭐ 优先级：中
- [ ] API 代码生成
- [ ] RPC 代码生成
- [ ] Model 代码生成
- [ ] Docker 文件生成
- [ ] Kubernetes YAML 生成

#### 5.2 API 定义 ⭐ 优先级：中
- [ ] .api 文件语法
- [ ] 类型定义
- [ ] 路由定义
- [ ] 文档生成

---

### 📊 第六阶段：可观测性

#### 6.1 指标监控 (Metrics) ⭐ 优先级：中
- [ ] Prometheus 集成
- [ ] 自定义指标
- [ ] 指标收集

#### 6.2 链路追踪 (Tracing) ⭐ 优先级：中
- [ ] Jaeger 集成
- [ ] Zipkin 集成
- [ ] 链路上下文传递

#### 6.3 日志聚合 ⭐ 优先级：低
- [ ] 日志收集
- [ ] 日志分析

---

### 🧪 第七阶段：测试与部署

#### 7.1 单元测试 ⭐ 优先级：高
- [ ] Mock 工具
- [ ] 测试覆盖率

#### 7.2 集成测试 ⭐ 优先级：中
- [ ] API 测试
- [ ] RPC 测试

#### 7.3 部署 ⭐ 优先级：低
- [ ] Docker 化
- [ ] Kubernetes 部署
- [ ] CI/CD 流程

## 📖 学习日志

### 2025-11-17

**学习内容：**
- ✅ 配置 VS Code 调试环境（安装 Delve）
- ✅ 创建 launch.json 调试配置
- ✅ 实现 color 包（终端颜色输出）
- ✅ 完善 logx Writer 测试用例
- ✅ 制定全面的学习计划（基于 go-zero 完整架构）

**学习要点：**
- **调试工具**：Delve 是 Go 的标准调试器
  - 支持断点、变量查看、单步执行
  - VS Code 集成良好
- **终端颜色**：使用 ANSI 转义序列
  - 前景色：`\033[3Xm`
  - 背景色：`\033[4Xm`
  - 重置：`\033[0m`

**代码示例：**
```go
// 使用颜色输出
colored := color.WithColor("错误", color.FgRed)
fmt.Println(colored) // 红色文本

// 带背景色和内边距
padded := color.WithColorPadding("成功", color.BgGreen)
fmt.Println(padded) // 绿色背景，带空格内边距
```

**思考与总结：**
- go-zero 的架构非常完整，涵盖了微服务的方方面面
- 学习需要有计划地推进，从核心包开始逐步深入
- 每个模块都有其独特的设计思想和工程实践

**下一步计划：**
- 完善 logx 的完整日志级别实现
- 实现日志轮转功能
- 开始 syncx 更多并发原语的实现

### 2025-11-16

**学习内容：**
- ✅ 初始化项目，创建 git 仓库并推送到 GitHub
- ✅ 完成 Logger 接口设计
- ✅ 实现延迟求值的 `Debugfn` 方法
- ✅ 为 Debugfn 添加详细的使用示例和注释
- ✅ 创建项目 README 和 LEARNING 文档
- ✅ 分离学习计划到独立文件

**学习要点：**
- **延迟求值**：通过 `func() any` 实现懒加载式日志记录
  - 只有在日志级别满足条件时才执行函数
  - 避免不必要的性能开销（如昂贵的计算、大对象格式化）
  - 适用场景：debug 日志、复杂的字符串拼接、序列化操作

**代码示例：**
```go
// 不使用延迟求值 - 即使 Debug 未启用也会执行
logger.Debug(fmt.Sprintf("data: %v", expensiveComputation()))

// 使用延迟求值 - 只在需要时执行
logger.Debugfn(func() any {
    return fmt.Sprintf("data: %v", expensiveComputation())
})
```

**思考与总结：**
- 性能优化的核心：避免不必要的计算
- 接口设计要考虑使用场景和性能权衡
- go-zero 在日志系统中大量使用这种模式
- 文档结构很重要：README 面向用户，LEARNING 记录学习过程

**下一步计划：**
- 实现完整的日志级别控制
- 添加结构化日志支持
- 学习 Writer 的不同实现

### 2025-11-15

**学习内容：**
- ✅ 学习 go-zero 日志系统源码
- ✅ 设计 logx 模块架构
- ✅ 实现 AtomicBool 原子操作

**学习要点：**
- 原子操作的使用场景
- 无锁编程的优势
- sync/atomic 包的使用

**思考与总结：**
- 原子操作适合简单的状态管理
- 比互斥锁性能更高
- 需要注意内存对齐问题

---

## 💡 核心知识点总结

### 日志系统设计

**延迟求值模式**
- 原理：使用闭包延迟执行昂贵操作
- 优势：提升性能，减少不必要的计算
- 实现：通过 `func() any` 类型参数

**日志级别**
- Debug: 调试信息
- Info: 一般信息
- Warning: 警告信息
- Error: 错误信息
- Fatal: 致命错误

**结构化日志**
- 便于日志分析和查询
- 支持 JSON 格式输出
- 易于集成日志收集系统

### 并发编程

**原子操作**
- 无锁的线程安全操作
- 适用于简单类型的并发访问
- 比互斥锁性能更高

**并发原语**
- Mutex: 互斥锁
- RWMutex: 读写锁
- WaitGroup: 等待组
- Once: 单次执行

**Channel 模式**
- 生产者-消费者
- 扇入扇出
- 超时控制

## 📚 参考资料

### 官方文档
- [go-zero 官方文档](https://go-zero.dev/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [Go 语言官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)

### 推荐阅读
- 《Go 语言设计与实现》
- 《Go 语言高级编程》
- go-zero 源码阅读

### 相关文章
- [go-zero 日志系统设计](https://go-zero.dev/docs/tutorials)
- [Go 并发编程实战](https://go.dev/blog/)

---

**记录时间线：每天坚持学习，记录成长的每一步 🚀**
