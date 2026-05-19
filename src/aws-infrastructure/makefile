.DEFAULT_GOAL = build

-include env.mk

build:
	go build -o provision *.go
	./provision

teardown:
	go build -o provision *.go
	./provision -teardown
