build: format
	go build -o bin/ ./...

format: generate
	gofumpt -l -w .
	gci write .

generate: restore
	go generate ./...

restore:
	go get -v ./...
	go mod tidy -v

