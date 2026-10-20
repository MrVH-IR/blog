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
JWT_SECRET=replace_this_with_a_long_random_secret(32Bytes)
JWT_TTL=15m
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

### DB
Tables
```sql
images,
posts,
users,
tags
```
Queries Can be found in every repository in internal/*/repository.go
Request Types can be found in internal/*/model.go

### For Testing Use cURL or Postman