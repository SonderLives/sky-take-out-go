# sky-take-out-go —— 黑马苍穹外卖 Go 版本

⎛⎝≥⏝⏝≤⎛⎝⎛⎝≥⏝⏝≤⎛⎝⎛⎝≥⏝⏝≤⎛⎝


[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-Web%20Framework-00A86B)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)

`sky-take-out-go` 使用 **Golang 后端脚手架 [GoFlow](https://github.com/s8sg/goflow)来完成快速开发**，

- **面向高并发 API 服务**：支持连接池、超时、重试、限流，完整的日志链路。
- **清晰分层架构**：`Handler -> Service -> Repository`，易于理解和扩展。
- **完善的中间件**：包含 RequestID、Logger、Recovery、CORS、I18n、Auth、Rate Limit。
- **高可扩展性**：通过 `ServiceContext` 实现统一依赖注入，方便扩展。
- **友好的可观测性**：使用 Zap 结构化日志，支持数据库和 Redis 的日志钩子。

## 核心技术

- **Gin Web API 框架**：快速构建高效的 Web API。
- **GORM + MySQL**：简化数据库操作，支持关系型数据存取。
- **Redis 缓存与基于 Redis Lua 脚本的路由限流**：确保高并发请求下的稳定性。
- **JWT 认证**：用于用户和管理员的身份认证。
- **Watermill MQ**：支持基于 Redis Streams 和 MySQL 后端的消息队列。
- **I18n 校验与错误消息**：支持多语言校验和错误提示。
- **优雅停机**：保证服务能够在停机时优雅退出，减少中断。

## 架构

```text
.
├── cmd/
│   └── server/         # 服务启动入口
├── config/             # 配置文件
├── internal/
│   ├── config/         # 配置管理
│   ├── database/       # 数据库操作
│   ├── handler/        # 控制器层（路由处理）
│   ├── middleware/     # 中间件（如认证、日志等）
│   ├── model/          # 数据模型
│   ├── mq/             # 消息队列（如 Watermill）
│   ├── pkg/            # 工具包（如 errcode、response、validator、ratelimit）
│   ├── repository/     # 数据访问层
│   ├── router/         # 路由管理
│   ├── service/        # 业务逻辑层
│   └── svc/            # ServiceContext 依赖注入
├── migration/          # 数据库迁移
├── go.mod              # Go 依赖管理
└── Makefile            # 构建与管理脚本
```

### Run

```bash
git clone https://github.com/your-username/goflow.git
cd goflow

go mod tidy
make run
```

Health check:

```bash
curl http://localhost:8080/health
```

## License

MIT License. See [LICENSE](./LICENSE).
