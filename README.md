# mini-zero

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个基于 [go-zero](https://github.com/zeromicro/go-zero) 的 Go 微服务框架学习项目。通过逐步实现 go-zero 的核心功能，深入理解微服务架构和 Go 语言最佳实践。

## 📚 项目简介

mini-zero 是一个学习型项目，旨在通过实现 go-zero 框架的核心组件来深入理解微服务框架的设计思想。项目采用渐进式开发，每个模块都经过仔细的设计和测试。

**学习目标：**
- 🔍 深入理解 go-zero 的核心设计理念
- 💡 掌握 Go 微服务开发的最佳实践
- 🛠️ 从零实现框架核心组件
- 📝 积累生产级代码编写经验

> 📖 详细的学习计划和日志请查看 [LEARNING.md](LEARNING.md)

## ✨ 已实现功能

### 日志系统 (logx)

提供高性能、易用的日志功能，支持延迟求值优化。

**核心特性：**
- ✅ Logger 接口设计
- ✅ 灵活的配置管理
- ✅ 延迟求值 (Debugfn) - 避免不必要的性能开销
- ✅ 多种日志级别 (Debug, Info, Warning, Error)

### 并发控制 (syncx)

提供高性能的并发原语和工具。

**核心特性：**
- ✅ AtomicBool - 原子布尔操作

## 🏗️ 项目结构

```
mini-zero/
├── core/
│   ├── logx/              # 日志系统
│   │   ├── config.go          # 配置定义
│   │   ├── fields.go          # 日志字段
│   │   ├── logger.go          # Logger 接口
│   │   ├── logs.go            # 日志实现
│   │   ├── logwriter.go       # 日志写入器
│   │   ├── vars.go            # 全局变量
│   │   ├── writer.go          # Writer 接口
│   │   └── *_test.go          # 单元测试
│   └── syncx/             # 并发控制
│       ├── atomicbool.go      # 原子布尔值
│       └── atomicbool_test.go # 单元测试
├── go.mod
├── README.md              # 项目说明
└── LEARNING.md            # 学习计划与日志
```

## 🚀 快速开始

### 前置要求

- Go 1.20 或更高版本

### 安装依赖

```bash
go mod download
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./core/logx
go test ./core/syncx

# 查看测试覆盖率
go test -cover ./...
```

### 使用示例

#### 日志系统

```go
package main

import "github.com/YunFy26/mini-zero/core/logx"

func main() {
    // 普通日志
    logx.Debug("simple debug message")
    logx.Info("application started")
    
    // 格式化日志
    logx.Debugf("user: %s, id: %d", "Alice", 123)
    
    // 延迟求值（推荐用于昂贵操作）
    logx.Debugfn(func() any {
        // 只在 Debug 级别启用时才执行
        return fmt.Sprintf("data: %v", computeExpensiveData())
    })
}
```

#### 原子操作

```go
package main

import "github.com/YunFy26/mini-zero/core/syncx"

func main() {
    // 创建原子布尔值
    flag := syncx.NewAtomicBool()
    
    // 设置值
    flag.Set(true)
    
    // 获取值
    if flag.True() {
        // do something
    }
    
    // 比较并交换
    flag.CompareAndSwap(true, false)
}
```

## 📚 参考资料

- [go-zero 官方文档](https://go-zero.dev/)
- [go-zero GitHub](https://github.com/zeromicro/go-zero)
- [学习计划与日志](LEARNING.md)

## 🤝 贡献

这是个人学习项目，欢迎提出建议和意见！

## 📄 许可证

MIT License

## 📬 联系方式

- GitHub: [@YunFy26](https://github.com/YunFy26)
- 项目地址: [mini-zero](https://github.com/YunFy26/mini-zero)

---

⭐ 如果这个项目对你有帮助，欢迎 Star！
