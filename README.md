# NebulaVPN

**Production-ready anti-DPI VPN service** with native clients for macOS, Windows, iOS, and Android.

## Architecture

```
┌──────────────────────────────────────────────────┐
│              CLIENTS                              │
│  macOS  │  Windows  │  iOS  │  Android           │
│  Swift  │  WinUI 3  │ Swift │  Kotlin            │
└─────────┴───────────┴───────┴───────────────────┘
                      │
            ┌─────────┴─────────┐
            │  Shared Rust Core  │
            │  (Transport/Crypto)│
            └───────────────────┘
                      │
              ┌───────┴───────┐
              │   CDN Relay   │
              │  (Cloudflare) │
              └───────────────┘
                      │
              ┌───────┴───────┐
              │  API Gateway  │
              │  (Go + Fiber) │
              └───────────────┘
                      │
              ┌───────┴───────┐
              │  Core Services │
              │  Auth/Device/  │
              │  Billing/Node  │
              └───────────────┘
                      │
              ┌───────┴───────┐
              │  Data Plane   │
              │  WG/VLESS/H2/ │
              │  TUIC         │
              └───────────────┘
```

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Backend | Go 1.22, Fiber, gRPC |
| Database | PostgreSQL 16, pgx |
| Cache | Redis 7 |
| Client Core | Rust (cross-platform) |
| Desktop | SwiftUI (macOS), WinUI 3 (Windows) |
| Mobile | SwiftUI (iOS), Kotlin + Compose (Android) |
| VPN Transport | WireGuard + VLESS+Reality + Hysteria2 + TUIC |
| Monitoring | Prometheus + Grafana + Loki |
| Deployment | Docker, Kubernetes, Terraform, Ansible |
| CI/CD | GitHub Actions |
| Billing | Stripe, Apple IAP, Google IAP |

## Protocols

| Protocol | Transport | Anti-DPI | Use Case |
|----------|-----------|----------|----------|
| WireGuard | UDP | Low | Primary (fastest) |
| VLESS+Reality | TCP+TLS | High | DPI-heavy regions |
| Hysteria2 | UDP (BBR) | Medium | Lossy networks |
| TUIC | TCP+TLS | Medium | Stable fallback |

## Quick Start

```bash
# Clone & setup
git clone https://github.com/nebula-vpn/nebula-vpn && cd nebula-vpn

# Start infrastructure
docker compose -f infrastructure/docker/docker-compose.yml up -d

# Run database migrations
cd backend && make migrate-up

# Start API server
cd backend && go run ./cmd/api-gateway

# Server runs at http://localhost:8080
```

## Project Structure

```
nebula-vpn/
├── backend/           # Go backend monorepo
├── clients/           # Native client apps
├── web/               # Admin & user panels
├── infrastructure/    # Docker, K8s, Terraform
├── docs/              # Architecture & API docs
└── scripts/           # Development scripts
```

## API

REST API available at `https://api.nebula-vpn.io/v1`

- `POST /auth/register` — Create account
- `POST /auth/login` — Get JWT tokens
- `GET /servers` — List VPN servers
- `GET /devices` — List devices
- `POST /devices` — Register device
- `GET /servers/recommend` — Smart server selection

## License

MIT
