module github.com/arbhalerao/cohereDB/db_manager

go 1.23.0

replace github.com/arbhalerao/cohereDB/db => ../db

replace github.com/arbhalerao/cohereDB/utils => ../utils

replace github.com/arbhalerao/cohereDB/config => ../config

replace github.com/arbhalerao/cohereDB/pb => ../pb

require (
	github.com/arbhalerao/cohereDB/pb v0.0.0-00010101000000-000000000000
	github.com/arbhalerao/cohereDB/utils v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/prometheus/client_golang v1.23.2
	github.com/rs/zerolog v1.33.0
	google.golang.org/grpc v1.70.0
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.8 // indirect
)
