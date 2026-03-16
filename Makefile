.PHONY: build run test clean

build:
	go build -o bin/smart-city-orchestrator .

run: build
	./bin/smart-city-orchestrator

test:
	go test ./...

clean:
	rm -rf bin/
