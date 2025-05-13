# cohereDB

A distributed key-value database with consistent hashing and dynamic server management.

## Architecture

```
┌─────────────────┐     ┌───────────────────────────┐
│   DB Manager    │     │        Data Servers       │
│  (Coordinator)  │────▶│  ┌─────┐ ┌─────┐ ┌─────┐  │
│                 │     │  │Node1│ │Node2│ │Node3│  │
│   gRPC: 9090    │     │  └─────┘ └─────┘ └─────┘  │
│   HTTP: 8090    │     └───────────────────────────┘
└─────────────────┘
        ▲
        │ gRPC
        │
┌─────────────┐
│   Client    │
│ (CLI Tool)  │
└─────────────┘
```

### Components

**DB Manager**  
- Routes requests using consistent hashing
- Manages data server registry and health monitoring  
- Single coordinator for the cluster

**DB Servers**  
- Store data using BadgerDB
- Register with manager on startup
- Handle CRUD operations

**Client**  
- CLI tool for database operations
- Connects to manager via gRPC

## Key Flows

### 1. Write Operation
```
Client ──SET──▶ Manager ──hash──▶ Target Server ──store──▶ BadgerDB
       "user:1"         GetNode()                 data
```

**Steps:**
1. Client sends `SetKey("user:1", "Alice")` to Manager
2. Manager uses consistent hashing: `GetNode("user:1")` → `server-2`  
3. Manager forwards request to target server
4. Server stores in BadgerDB
5. Response flows back to client

### 2. Server Registration
```
New Server ──register──▶ Manager ──update──▶ Hash Ring ──start──▶ Health Checks
 startup      HTTP                add node              monitor
```

**Steps:**
1. Server starts and calls `POST /register` with details
2. Manager adds server to registry and hash ring
3. Manager begins health monitoring
4. Server starts receiving requests

### 3. Consistent Hashing
```
Key ──hash──▶ Position on Ring ──lookup──▶ Assigned Server
    "user:1"      hash: 250           └──▶ server-2
```

## Current State

### ✅ **Implemented**
- Distributed storage with consistent hashing
- Service discovery and health monitoring
- Basic CRUD operations (GET, SET, DELETE)

### ❌ **Missing** 
- Data replication (single copy per key)
- Manager redundancy (single point of failure)
- Dynamic scaling with data migration
