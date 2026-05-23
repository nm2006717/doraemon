# Doraemon - Personal Knowledge Base MCP Server

A personal knowledge base inspired by Doraemon's 4D pocket. Store anything, find it instantly.

## Project Overview

- **Language**: Go
- **Storage**: SQLite with FTS5 full-text search
- **File Storage**: Local filesystem (`data/files/`), no MinIO
- **Interface**: MCP Server (both store and retrieve via MCP)
- **Deployment**: Local binary + Docker
- **VCS**: Git

## MCP Tools

| Tool | Purpose |
|------|---------|
| `store` | Store notes/files/links/images with auto-classification and tagging |
| `search` | Search by keyword, tag, or content |
| `list` | Browse all entries or filter by category |
| `recommend` | Suggest relevant entries based on current conversation context |

## Supported Content Types

- Plain text (ideas, notes, code snippets)
- Files (PDF, images, documents)
- Links/bookmarks
- Images

## Architecture

```
doraemon/
├── cmd/              # Entry point
├── internal/
│   ├── mcp/          # MCP server protocol handling
│   ├── store/        # Storage layer (SQLite + file system)
│   ├── classify/     # Auto-classification logic
│   └── search/       # Search and recommendation engine
├── data/
│   ├── doraemon.db   # SQLite database
│   └── files/        # Stored files
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── CLAUDE.md
```

## Development

```bash
# Run locally
go run ./cmd/doraemon

# Build
go build -o doraemon ./cmd/doraemon

# Docker
docker compose up
```

## Design Decisions

- No separate CLI or Web UI needed — all interaction happens through MCP
- Auto-classification uses keyword extraction and simple rules (no external AI dependency for tagging)
- SQLite FTS5 for full-text search — lightweight, no extra service needed
- Files stored on local filesystem, metadata in SQLite
