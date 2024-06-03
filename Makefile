default: build

clean:
	rm -rf out

test:
	go test ./...

build:
	mkdir -p out
	go build -o out/example ./cmd/painter/main.go
