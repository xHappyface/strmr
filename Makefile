clean:
	rm -rf vendor
	rm -f go.sum
	go get ./...
	go mod vendor

obs: clean
	go run cmd/obs/main.go --task="New Cool Task"

run: clean
	go run main.go