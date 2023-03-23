.PHONY: update test

update:
	go run ./pkg/...

test:
	go test ./...
