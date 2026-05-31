.DEFAULT_GOAL = run

-include env.mk
-include dsql.mk

compile:
	go build -ldflags="-w -s -X main.GIT_COMMIT_HASH=`git rev-parse HEAD`" -o schema *.go

run: compile
	./schema -id=999
	./schema --list-tables

rollback: compile
	./schema -id=3 --rollback

flatten: compile
	./schema -id=0 --rollback

ls: compile
	./schema --list-tables

environments: compile
	./schema --environments

data: compile
	./schema -id=49 --rollback
	./schema -id=50