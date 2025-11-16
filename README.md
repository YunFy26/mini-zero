# mini-zero

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个基于 [go-zero](https://github.com/zeromicro/go-zero) 的 Go 微服务框架学习项目。通过逐步实现 go-zero 的核心功能，深入理解微服务架构和 Go 语言最佳实践。

## 📚 项目简介

mini-zero 是我学习 go-zero 框架的实践项目，目标是：

- 🔍 深入理解 go-zero 的核心设计理念
- 💡 掌握 Go 微服务开发的最佳实践
- 🛠️ 从零实现框架核心组件
- 📝 记录学习过程和技术心得

## 🎯 学习计划

### 第一阶段：基础组件（进行中）

- [x] **日志系统 (logx)**
  - [x] Logger 接口设计
  - [x] 配置管理
  - [x] 延迟求值优化
  - [ ] 多种 Writer 实现
  - [ ] 日志级别控制
  - [ ] 结构化日志

- [ ] **并发控制 (syncx)**
  - [x] AtomicBool 实现
  - [ ] 限流器 (Rate Limiter)
  - [ ] 熔断器 (Circuit Breaker)
  - [ ] 线程池

### 第二阶段：网络通信

- [ ] **HTTP 服务器**
  - [ ] 路由设计
  - [ ] 中间件机制
  - [ ] 请求响应处理

- [ ] **RPC 框架**
  - [ ] 服务注册与发现
  - [ ] 负载均衡
  - [ ] 超时控制

### 第三阶段：高级特性

- [ ] **缓存系统**
- [ ] **数据库集成**
- [ ] **链路追踪**
- [ ] **服务监控**

## 📖 学习日志

### 2025-11-16
- ✅ 初始化项目，创建 git 仓库
- ✅ 完成 Logger 接口设计
- ✅ 实现延迟求值的 `Debugfn` 方法，添加使用示例
- ✅ 学习要点：
  - 延迟求值可以避免不必要的性能开销
  - 通过 `func() any` 实现懒加载式日志记录

### 2025-11-15
- ✅ 学习 go-zero 日志系统源码
- ✅ 设计 logx 模块架构
- ✅ 实现 AtomicBool 原子操作

## 🏗️ 项目结构

```
mini-zero/
├── core/
│   ├── logx/          # 日志系统
│   │   ├── config.go      # 配置定义
│   │   ├── logger.go      # Logger 接口
│   │   ├── writer.go      # 日志写入器
│   │   └── ...
│   └── syncx/         # 并发控制
│       ├── atomicbool.go  # 原子布尔值
│       └── ...
├── go.mod
└── README.md
```

## 🚀 快速开始

### 安装依赖

```bash
go mod download
```

### 运行测试

```bash
go test ./...
```

### 使用示例

```go
package main

import "github.com/YunFy26/mini-zero/core/logx"

func main() {
    logger := logx.NewLogger()
    
    // 普通日志
    logger.Debug("simple debug message")
    
    // 格式化日志
    logger.Debugf("user: %s, id: %d", "Alice", 123)
    
    // 延迟求值（只在 Debug 级别启用时才执行）
    logger.Debugfn(func() any {
        // 昂贵的计算只在需要时执行
        return computeExpensiveData()
    })
}
```

## 📚 参考资料

- [go-zero 官方文档](https://go-zero.dev/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [Go 语言官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)

## 💡 核心知识点

### 日志系统设计

- **延迟求值**：通过闭包实现懒加载，避免不必要的性能开销
- **结构化日志**：便于日志分析和查询
- **日志级别**：Debug, Info, Warning, Error, Fatal

### 并发编程

- **原子操作**：无锁的线程安全操作
- **并发原语**：Mutex, RWMutex, WaitGroup
- **channel 模式**：生产者-消费者、扇入扇出

## 🤝 贡献

这是个人学习项目，欢迎提出建议和意见！

## 📄 许可证

MIT License

## 📬 联系方式

- GitHub: [@YunFy26](https://github.com/YunFy26)
- 项目地址: [mini-zero](https://github.com/YunFy26/mini-zero)

---

**学习笔记**：每天坚持学习一点，记录成长的每一步 🚀
