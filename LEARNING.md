# mini-zero 学习计划与日志

## 🎯 学习目标

通过逐步实现 go-zero 的核心功能，深入理解微服务框架设计和 Go 语言最佳实践。

## 📅 学习计划

### 第一阶段：基础组件（进行中）

#### 日志系统 (logx)
- [x] Logger 接口设计
- [x] 配置管理
- [x] 延迟求值优化 (Debugfn)
- [ ] 多种 Writer 实现
- [ ] 日志级别控制
- [ ] 结构化日志
- [ ] 日志轮转
- [ ] 异步写入优化

#### 并发控制 (syncx)
- [x] AtomicBool 实现
- [ ] AtomicDuration
- [ ] 限流器 (Rate Limiter)
- [ ] 熔断器 (Circuit Breaker)
- [ ] 线程池 (Pool)
- [ ] 资源管理器

### 第二阶段：网络通信

#### HTTP 服务器
- [ ] 路由设计与实现
- [ ] 中间件机制
- [ ] 请求响应处理
- [ ] 参数绑定与验证
- [ ] 错误处理
- [ ] 超时控制

#### RPC 框架
- [ ] 服务注册与发现
- [ ] 负载均衡策略
- [ ] 超时与重试
- [ ] 服务治理

### 第三阶段：高级特性

#### 缓存系统
- [ ] 本地缓存
- [ ] Redis 集成
- [ ] 缓存策略

#### 数据库集成
- [ ] MySQL/PostgreSQL 支持
- [ ] 连接池管理
- [ ] ORM 集成

#### 可观测性
- [ ] 链路追踪 (Tracing)
- [ ] 服务监控 (Metrics)
- [ ] 健康检查

## 📖 学习日志

### 2025-11-16

**学习内容：**
- ✅ 初始化项目，创建 git 仓库并推送到 GitHub
- ✅ 完成 Logger 接口设计
- ✅ 实现延迟求值的 `Debugfn` 方法
- ✅ 为 Debugfn 添加详细的使用示例和注释
- ✅ 创建项目 README 文档

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
