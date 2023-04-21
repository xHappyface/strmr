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

db:
	sqlite3 strmr.db < sql/schema.sql

auth-yt:
	sudo apt install python3-pip
	pip install --upgrade google-api-python-client google-auth-httplib2 google-auth-oauthlib
	source scripts/token.sh && python3 scripts/auth.py