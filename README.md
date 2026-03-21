# Dorgu Platform

Real-time web dashboard for visualizing and managing Kubernetes clusters via the ClusterPersona CRD.

## Overview

Dorgu Platform provides a web-based UI for viewing cluster state powered by the Dorgu Operator's ClusterPersona CRD. It includes:

- **Cluster List View**: See all managed clusters with phase, environment, nodes, and stack status
- **Cluster Detail View**: Deep dive into cluster resources, installed components, and node health
- **Real-time Updates**: WebSocket-based live updates as the operator reconciles cluster state
- **Embeddable**: Can be embedded into the `dorgu` CLI or run as a standalone service

## Architecture

- **Backend**: Go (embeddable package + standalone server)
- **Frontend**: React + TypeScript + Vite + Tailwind CSS + shadcn/ui
- **State Management**: React Query (TanStack Query)
- **Real-time**: WebSocket (client-go Informer → broadcast)
- **Deployment**: Local (`dorgu platform serve`) or cluster-based (future)

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- kubectl configured with access to a cluster running Dorgu Operator
- ClusterPersona CRDs installed

### Run Locally (Development)

```bash
# Backend
make run

# Frontend (separate terminal)
cd web
npm install
npm run dev
```

Navigate to `http://localhost:8080`

### Build

```bash
# Build backend with embedded frontend
make build

# Run the binary
./bin/dorgu-platform serve
```

### Use via dorgu CLI (Recommended)

```bash
# Start the platform
dorgu platform serve

# With custom port
dorgu platform serve --port 3000

# With specific kubeconfig/context
dorgu platform serve --kubeconfig ~/.kube/custom-config --context prod-cluster
```

## Project Structure

```
dorgu-platform/
├── cmd/
│   └── server/          # Standalone server binary
├── pkg/
│   ├── platform/        # Public API for embedding in dorgu CLI
│   ├── server/          # HTTP server and static file serving
│   ├── api/             # HTTP handlers
│   ├── websocket/       # WebSocket hub and client
│   ├── watcher/         # K8s client-go Informer for ClusterPersona
│   └── models/          # Data models
├── web/                 # React frontend
│   ├── src/
│   ├── public/
│   └── package.json
├── Makefile
└── go.mod
```

## Development

### Make Targets

```bash
make help          # Show all available targets
make build         # Build backend + frontend
make run           # Run backend in development mode
make test          # Run Go tests
make lint          # Run linters
make clean         # Clean build artifacts
make frontend      # Build frontend only
```

## Testing

```bash
# Backend tests
make test

# Frontend tests
cd web
npm run test
```

## API Endpoints

- `GET /api/clusters` — List all ClusterPersona resources
- `GET /api/clusters/:name` — Get single ClusterPersona
- `GET /ws` — WebSocket endpoint for real-time updates

## Embedding

The platform can be embedded into any Go application via the `pkg/platform` package:

```go
import "github.com/dorgu-ai/dorgu-platform/pkg/platform"

config := platform.Config{
    Port:       "8080",
    Kubeconfig: "",       // defaults to $KUBECONFIG or ~/.kube/config
    Context:    "",       // defaults to current context
}

srv, err := platform.NewServer(config)
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()
if err := srv.Start(ctx); err != nil {
    log.Fatal(err)
}
```

## Phase 1 MVP Scope

- Real-time ClusterPersona visualization
- Cluster list and detail views
- WebSocket live updates
- Local development mode
- Embeddable in dorgu CLI

## Future Phases

- **Phase 2**: Multi-cluster support, cluster deployment (Helm chart)
- **Phase 3**: Incident tracking, irregularity flagging
- **Phase 4**: AI agent decision log visualization
- **Phase 5**: Autonomous reconciliation modes

## License

Apache 2.0 (same as dorgu/dorgu-operator)

## Contributing

See the main [dorgu repository](https://github.com/dorgu-ai/dorgu) for contribution guidelines.
