# A Simple Blog with Golang

A straightforward blog built with Go. This project is an API only and has no frontend.

## Requirements

- Go 1.26
- PostgreSQL

## Configuration

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/blog?sslmode=disable
HTTP_ADDR=localhost:8000
```

## Running the Project

```bash
go mod tidy
go run ./cmd/api
```

## Testing

Since there is no frontend, you are welcome to test the API with `curl`.

For example, to check the database connection:

```bash
curl -X GET "localhost:8000/health/db"
```