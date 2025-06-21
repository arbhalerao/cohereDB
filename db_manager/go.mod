module github.com/arbhalerao/cohereDB/db_manager

go 1.22.0

replace github.com/arbhalerao/cohereDB/db => ../db

replace github.com/arbhalerao/cohereDB/utils => ../utils

replace github.com/arbhalerao/cohereDB/config => ../config

replace github.com/arbhalerao/cohereDB/pb => ../pb

require (
	github.com/arbhalerao/cohereDB/pb v0.0.0-00010101000000-000000000000
	github.com/arbhalerao/cohereDB/utils v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.70.0
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
