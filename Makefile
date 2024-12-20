run:
	@go run ./cmd/app/main.go

bench:
	@go test ./cmd/app/... -bench=. -benchmem

test:
	@go test ./cmd/app/... -v

profile:
	@go run ./cmd/app/main.go -profile
	@go tool pprof -http 127.0.0.1:8080 cpu_profile.prof
