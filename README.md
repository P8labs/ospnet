# ospnet

> A distributed network for deploying apps on independent machines.

ospnet is a lightweight, decentralized platform that allows developers to deploy and run applications across a network of independent nodes (Raspberry Pis, home servers, and personal machines).

It is designed for **side projects, demos, and experimentation**, not production-grade workloads.

## Why ospnet?

Deploying small projects today is unnecessarily complex:

- PaaS platforms introduce vendor lock-in
- Serverless platforms sleep or limit usage
- Cloud infrastructure is overkill for demos

ospnet provides a simple alternative:

- Bring your own hardware
- Run Docker containers directly
- Share compute across a trusted network

## Core Idea

ospnet connects multiple independent machines into a single logical network.

Each node contributes resources (CPU, RAM, storage), and applications are scheduled across available nodes.

A central controller manages:

- Node registration
- Health monitoring
- Deployment scheduling
- Traffic routing

## Architecture

### 1. Master Node

- API server
- Scheduler
- Reverse proxy
- Network coordinator

### 2. Worker Nodes

- Run ospnet agent
- Connect via secure mesh network
- Execute Docker containers
- Report health and resource usage

### 3. Networking

- Nodes communicate over a private mesh network
- Public traffic is routed through a central entry point

## Features (MVP)

- Node registration and discovery
- Docker-based deployments
- Basic scheduling (1 app → 1 node)
- Health checks and restart on failure
- Public URL via subdomain routing

## Non-Goals

ospnet is **not** intended to be:

- A production-grade cloud platform
- A Kubernetes replacement
- A high-availability system
- A commercial hosting solution

## Roadmap

- [x] Node agent (Go)
- [ ] Master API
- [ ] Basic scheduler
- [ ] Reverse proxy integration
- [ ] CLI tool

## Agent Daemon (MVP)

The OSPNet node agent is a long-running daemon that:

- onboards against master using a token
- joins Tailscale using returned auth key
- exposes local node API (`/health`, `/containers`, `/containers/run`, `/containers/stop`)
- persists container state in SQLite
- runs heartbeat and reconciliation loops

### Build

```bash
go build -o ospnet-agent ./cmd/agent
```

### Environment

```bash
OSPNET_MASTER_URL=http://100.64.0.10:8080
OSPNET_TOKEN_PATH=/etc/ospnet/token
OSPNET_CONFIG_PATH=/etc/ospnet/config.json
OSPNET_DB_PATH=/var/lib/ospnet/agent.db
OSPNET_AGENT_PORT=9000
```

Optional metadata:

```bash
OSPNET_NODE_NAME=pi-01
OSPNET_NODE_REGION=home-lab
OSPNET_NODE_TYPE=raspi
```

### Installer

```bash
curl -sSL https://ospnet.run/install.sh | bash
```

Reference installer and unit files are included in `scripts/install.sh` and `scripts/ospnet-agent.service`.

## License

TBH
