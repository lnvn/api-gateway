.DEFAULT_GOAL := run

# PHONY keep make from getting confused, if you ever create  a directory in project with the same name as a target
fmt:
	go fmt ./...
.PHONY: fmt

# build: fmt  -> run fmt before build
build: fmt
	go build .
.PHONY: build

run: fmt
	go run .
.PHONY: run
