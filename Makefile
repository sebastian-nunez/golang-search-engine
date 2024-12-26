build:
	@go build -o ./bin/golang-search-engine ./cmd/main.go

run: build
	@./bin/golang-search-engine

test:
	@go test -v ./...

coverage:
	@go test -cover ./...

coverage_ui:
	go test -coverprofile=c.out ./... && go tool cover -html=c.out