run:
	@go run ./cmd/app/main.go

bench:
	@go test ./cmd/profiler/... -bench=. -benchmem

test:
	@go test ./cmd/profiler/... -v

profile:
	@go run ./cmd/profiler/profiler.go
	@go tool pprof -http 127.0.0.1:8080 cpu_profile.prof
