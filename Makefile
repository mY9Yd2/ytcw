all: ytcwd migrate

ytcwd:
	go build -o build/ytcwd ./cmd/daemon

migrate:
	go build -o build/ytcw-migrate ./cmd/migrate

clean:
	go clean
	rm -rf ./build/
