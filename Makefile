.PHONY: run test

run:
	LOCAL_DEV=true go run server.go

test:
	go test ./... -coverprofile cover.out
