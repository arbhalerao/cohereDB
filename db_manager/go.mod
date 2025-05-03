module github.com/arbha1erao/cohereDB/db_manager

go 1.22.0

replace github.com/arbha1erao/cohereDB/db => ../db

replace github.com/arbha1erao/cohereDB/utils => ../utils

replace github.com/arbha1erao/cohereDB/config => ../config

require (
	github.com/arbha1erao/cohereDB/utils v0.0.0-20250219115734-b825cb0f3e93
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
)
