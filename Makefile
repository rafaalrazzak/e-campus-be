db_migrate:
	go run cmd/bun/main.go db migrate

test:
	go test ./...

fmt:
	gofmt -w -s ./
	goimports -w  -local ecampus-be ./
