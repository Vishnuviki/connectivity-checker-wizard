.PHONY: run test

run:
	go run server.go

test:
	go test ./... -coverprofile cover.out
