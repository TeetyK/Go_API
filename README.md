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

---

## SMTP Email for Password Resets

The password reset functionality has been updated to send emails using SMTP.

To use this feature, you need to:

1.  **Add the `gomail` dependency:**
    ```bash
    go get gopkg.in/gomail.v2
    ```

2.  **Tidy the go modules:**
    ```bash
    go mod tidy
    ```

3.  **Configure the SMTP server settings via environment variables:**
    - `SMTP_HOST`: The SMTP server hostname (e.g., `smtp.example.com`).
    - `SMTP_PORT`: The SMTP server port (e.g., `587`).
    - `SMTP_USER`: Your SMTP username.
    - `SMTP_PASS`: Your SMTP password or app password.

If these environment variables are not set, the application will fall back to printing the password reset link to the console.