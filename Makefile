.PHONY: build, run, test

build:
	go build -o ./bin/budget_manager.exe ./cmd/budget_manager/*.go

run: build
	./bin/budget_manager.exe

test:
	go test ./... -v