default: start

build:
	@cd web && yarn && yarn build

start: build
	@go run ./cmd/...

mods:
	@go mod download

test:
	@go test ./...