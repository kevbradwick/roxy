.PHONT: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...

.PHONY: build-image
build-image:
	docker build -t roxy:latest .