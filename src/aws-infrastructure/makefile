.DEFAULT_GOAL = build

-include env.mk

b:
	go build -o provision *.go

build:
	go build -o provision *.go
	./provision

status:
	go build -o provision *.go
	./provision -status

teardown:
	go build -o provision *.go
	./provision -teardown
