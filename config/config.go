package config

type Config struct {
	Database    DatabaseConfig `toml:"database"`
	Server      ServerConfig   `toml:"server"`
	PeerServers []ServerConfig `toml:"peer_servers"`
}

type DatabaseConfig struct {
	ShardCount int `toml:"shard_count"`
}

type ServerConfig struct {
	Shard         int    `toml:"shard"`
	Addr          string `toml:"addr"`
	ContainerAddr string `toml:"container_addr"`
}
