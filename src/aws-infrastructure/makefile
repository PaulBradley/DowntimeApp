.DEFAULT_GOAL = provision

-include env.mk

compile:
	go build -ldflags="-w -s -X main.GIT_COMMIT_HASH=`git rev-parse HEAD`" -o provision *.go

provision: compile
	./provision

status: compile
	./provision -status

teardown: compile
	./provision -teardown
