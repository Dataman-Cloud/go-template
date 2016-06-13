default: help

build:

run:

localbuild:
	go build .

localrun:
	./go-template

local: localbuild localrun


help:
	@echo "demo project for sry"
	@echo "================"
	@echo "build      - build with docker"
	@echo "run        - run with docker"
	@echo "localbuild - build binary"
