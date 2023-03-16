build:
	go build -o bin/kinda-store

run: build
	./bin/kinda-store

test:
	go test -v ./...