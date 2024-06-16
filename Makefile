run: build
	@./bin/app --config config.toml

build:
	@go build -o bin/app cmd/main.go

