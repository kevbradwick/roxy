.PHONT: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...