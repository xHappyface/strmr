module github.com/jnrprgmr/dog

go 1.19

replace (
    github.com/nicklaw5/helix/v2 => /home/jnrprgmr/Projects/helix
)

require (
	github.com/andreykaipov/goobs v0.12.0
	github.com/nicklaw5/helix/v2 v2.16.0
	github.com/spf13/cobra v1.6.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/mattn/go-sqlite3 v1.14.16
)

require (
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
