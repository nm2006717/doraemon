# Doraemon

[English](#english) | [中文](#中文)

---

## English

A personal knowledge base inspired by Doraemon's 4D pocket. Store anything, find it instantly.

### Features

- Knowledge entry management (notes, links, files, images)
- Full-text search (SQLite FTS5)
- Auto-classification and tag extraction
- Web admin panel (Vue 3 + Tailwind)
- MCP Server interface (interact directly via Claude Code and other MCP clients)

### Quick Start

#### Local

```bash
# Build frontend
cd web && npm install && npm run build && cd ..

# Build and run
go build -o doraemon ./cmd/doraemon
./doraemon
```

Visit `http://localhost:8080` after startup. First-time use will guide you through creating an admin account.

#### Docker

```bash
docker compose up -d
```

#### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DORAEMON_HTTP_PORT` | HTTP server port | 8080 |
| `DORAEMON_JWT_SECRET` | JWT signing key (optional, randomly generated on each startup if not set) | random |

### Architecture

```
doraemon/
├── cmd/doraemon/       # Entry point, starts HTTP + MCP
├── internal/
│   ├── api/            # REST API (auth, CRUD, search)
│   ├── mcp/            # MCP Server (stdio)
│   ├── store/          # SQLite storage layer
│   ├── classify/       # Auto-classification
│   └── search/         # Recommendation engine
├── web/                # Vue 3 frontend (embedded in binary)
│   ├── src/
│   └── embed.go
├── data/               # Runtime data (SQLite + files)
├── Dockerfile
└── docker-compose.yml
```

### MCP Tools

Available via MCP protocol (stdio) for use in Claude Code and other MCP clients:

| Tool | Purpose |
|------|---------|
| `store` | Store notes, links, and files |
| `search` | Full-text search |
| `list` | Browse entries |
| `recommend` | Suggest relevant entries based on conversation context |

### Web Admin Panel

- JWT authentication
- Dashboard with statistics overview
- Entry list with pagination, category filtering, and search
- Create / edit / delete entries

### Tech Stack

- Backend: Go, net/http, SQLite (FTS5), JWT
- Frontend: Vue 3, TypeScript, Tailwind CSS, Vite
- Deployment: Single binary (frontend embedded via go:embed), Docker

---

## 中文

个人知识库，灵感来自哆啦A梦的四维口袋。存什么都行，找什么都快。

### 功能

- 知识条目管理（笔记、链接、文件、图片分类）
- 全文搜索（SQLite FTS5）
- 自动分类和标签提取
- Web 管理面板（Vue 3 + Tailwind）
- MCP Server 接口（通过 Claude Code 等工具直接交互）

### 快速开始

#### 本地运行

```bash
# 构建前端
cd web && npm install && npm run build && cd ..

# 编译运行
go build -o doraemon ./cmd/doraemon
./doraemon
```

启动后访问 `http://localhost:8080`，首次使用会引导创建管理员账号。

#### Docker

```bash
docker compose up -d
```

#### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DORAEMON_HTTP_PORT` | HTTP 服务端口 | 8080 |
| `DORAEMON_JWT_SECRET` | JWT 签名密钥（可选，不设则每次启动随机生成） | 随机 |

### 架构

```
doraemon/
├── cmd/doraemon/       # 入口，启动 HTTP + MCP
├── internal/
│   ├── api/            # REST API（认证、CRUD、搜索）
│   ├── mcp/            # MCP Server（stdio）
│   ├── store/          # SQLite 存储层
│   ├── classify/       # 自动分类
│   └── search/         # 推荐引擎
├── web/                # Vue 3 前端（embed 进二进制）
│   ├── src/
│   └── embed.go
├── data/               # 运行时数据（SQLite + 文件）
├── Dockerfile
└── docker-compose.yml
```

### MCP 工具

通过 MCP 协议（stdio）提供以下工具，可在 Claude Code 中直接使用：

| 工具 | 用途 |
|------|------|
| `store` | 存储笔记/链接/文件 |
| `search` | 全文搜索 |
| `list` | 浏览条目 |
| `recommend` | 基于上下文推荐相关条目 |

### Web 管理面板

- 登录认证（JWT）
- Dashboard 统计概览
- 条目列表（分页、分类筛选、搜索）
- 新建/编辑/删除条目

### 技术栈

- 后端：Go、net/http、SQLite（FTS5）、JWT
- 前端：Vue 3、TypeScript、Tailwind CSS、Vite
- 部署：单二进制（前端 embed）、Docker
