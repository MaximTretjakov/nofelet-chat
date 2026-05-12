## build: Build an application
.PHONY: build_signaling
build_signaling:
	go build --ldflags="-s -w" -v cmd/signaling/main.go

## build: Run an application
.PHONY: run_signaling
run_signaling:
	go run cmd/signaling/main.go