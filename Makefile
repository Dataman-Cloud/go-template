default: help

build:

run:

localbuild:
	go build -o build/go-template .

localrun:
	./build/go-template

local: localbuild localrun

fmt:
	go fmt ./...


help:
	@echo "demo project for sry"
	@echo "================"
	@echo "build      - build with docker"
	@echo "run        - run with docker"
	@echo "localbuild - build binary"
