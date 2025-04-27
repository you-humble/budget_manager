.PHONY: build, run, test

build:
	go build -o budget_manager.exe ./cmd/budget_manager/*.go

run: build
	./budget_manager.exe

test:
	go test ./... -v