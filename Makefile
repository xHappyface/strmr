SHELL := /bin/bash

clean:
	rm -rf vendor
	rm -f go.sum
	go get ./...
	go mod vendor

obs: clean
	go run cmd/obs/main.go --task="Use HTML templates"

run: clean
	source scripts/token.sh && go run main.go