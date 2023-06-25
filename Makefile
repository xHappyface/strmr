SHELL := /bin/bash

deps:
	sudo apt install espeak
	sudo apt install python3-pip
	pip install --upgrade google-api-python-client google-auth-httplib2 google-auth-oauthlib

clean:
	rm -rf vendor
	rm -f go.sum
	go get ./...
	go mod vendor

obs: clean
	go run cmd/obs/main.go --task="Use HTML templates"

run: clean
	source scripts/token.sh && go run main.go

db-reset:
	rm -f strmr.db
	sqlite3 strmr.db < sql/schema.sql
	sqlite3 strmr.db < sql/data.sql

auth-yt:
	source scripts/token.sh && go run scripts/quickstart.go