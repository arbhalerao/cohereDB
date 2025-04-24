# cohereDB (WIP)

```
.
├── cli
│   ├── client
│   │   └── http_client.go        # HTTP client used for communication between the CLI and the server
│   ├── cmd
│   │   └── main.go               # Entry point for the CLI application
│   ├── commands
│   │   ├── delete.go             # Command handler for deleting key-value pairs
│   │   ├── get.go                # Command handler for retrieving key-value pairs
│   │   └── set.go                # Command handler for setting key-value pairs
│   ├── go.mod
│   └── go.sum
├── cohere
│   ├── cmd
│   │   └── main.go               # Entry point for the database server
│   ├── configs
│   │   ├── server0.toml          # Configuration file for server instance 0
│   │   ├── server1.toml          # Configuration file for server instance 1
│   │   └── server2.toml          # Configuration file for server instance 2
│   ├── db
│   │   └── db.go                 # Database logic (key-value store operations)
│   ├── utils
│   │   ├── logger.go             # Logging utilities
│   │   └── toml_loader.go        # TOML configuration loader
│   └── web
│       ├── config.go             # Web server configuration
│       ├── handlers.go           # HTTP request handlers
│       └── server.go             # Web server logic to expose the database as a service
│   ├── docker-compose.yaml       # Docker Compose setup for running multiple server instances
│   ├── Dockerfile                # Dockerfile for building the server container
│   ├── go.mod
│   ├── go.sum
└── README.md
```
