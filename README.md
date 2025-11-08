# Golang Example System


This is a minimal example of a Go REST API using Gin/GORM and PostgreSQL.


## Requirements
- Go 1.20+
- PostgreSQL running


## Setup
1. Copy `.env.example` to `.env` or set environment variables in your shell.
2. `go mod tidy`
3. `go run .` (or `go build` then run binary)


Endpoints:
- `GET /health` - health check
- `POST /users` - create user `{ "name": "Alice", "email": "a@b.com" }`
- `GET /users/:id` - get user by id


## Notes about common errors you mentioned
- `expected 'package', found 'EOF'` : make sure every `.go` file starts with `package <name>` and the file isn't truncated.
- `c.IndentadJSON undefined` : typo — Gin's method is `c.IndentedJSON(...)` (note the order and spelling).
- If you want to use a different PostgreSQL driver (pgx), update imports and DSN accordingly.