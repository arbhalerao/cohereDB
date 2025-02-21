module github.com/arbha1erao/cohereDB/db_server

go 1.22.0

replace github.com/arbha1erao/cohereDB/db => ../db

replace github.com/arbha1erao/cohereDB/utils => ../utils

replace github.com/arbha1erao/cohereDB/config => ../config

replace github.com/arbha1erao/cohereDB/pb => ../pb

require (
	github.com/arbha1erao/cohereDB/db v0.0.0-00010101000000-000000000000
	github.com/arbha1erao/cohereDB/pb v0.0.0-00010101000000-000000000000
	github.com/arbha1erao/cohereDB/utils v0.0.0-00010101000000-000000000000
	github.com/dgraph-io/badger v1.6.2
	google.golang.org/grpc v1.70.0
)

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgraph-io/badger/v4 v4.5.1 // indirect
	github.com/dgraph-io/ristretto v0.0.2 // indirect
	github.com/dgraph-io/ristretto/v2 v2.1.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/flatbuffers v24.12.23+incompatible // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
