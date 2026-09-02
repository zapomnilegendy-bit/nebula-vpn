# NebulaVPN Architecture

## Overview

NebulaVPN is a modern, anti-DPI resilient VPN service with native clients for macOS, Windows, iOS, and Android.

## Core Principles

1. **Transport Abstraction** — All protocols implement a common `Transport` interface
2. **Adaptive Protocol Selection** — ML-based scoring switches protocols based on real-time network conditions
3. **Layered Security** — Defense in depth from application to infrastructure
4. **Zero Trust Networking** — mTLS between services, E2E encryption where possible

## System Components

### Backend (Go)
- **API Gateway** — REST API entry point (Fiber)
- **gRPC Gateway** — Internal service mesh
- **Auth Service** — JWT + OAuth2 + Argon2id
- **Transport Router** — Adaptive protocol selection engine
- **Node Manager** — Server orchestration + health checks
- **Billing Service** — Stripe/Apple/Google IAP integration

### Data Plane
- **WireGuard** — Primary UDP transport
- **VLESS+Reality** — TCP with TLS camouflage
- **Hysteria2** — QUIC-based, optimized for lossy networks
- **TUIC** — TCP+TLS, stable fallback

### Client Core (Rust)
- Cross-platform shared library
- Transport implementation for all protocols
- Crypto operations (key generation, encryption)
- Config parsing and management

### Native Clients
- **macOS**: SwiftUI + NetworkExtension
- **Windows**: WinUI 3 + WireGuardNT
- **iOS**: SwiftUI + NEPacketTunnelProvider
- **Android**: Kotlin + VpnService

## Anti-DPI Stack

```
TLS Camouflage → Traffic Padding → Packet Fragmentation → Timing Randomization → JA3 Rotation
```

## Data Flow

```
App → Transport Layer → Obfuscation Layer → CDN Relay → VPN Server → Internet
```

## Database

PostgreSQL with the following core entities:
- users, devices, vpn_servers, plans, subscriptions, connections, promocodes

## Deployment

- Kubernetes (EKS/GKE) for control plane
- Bare metal / VMs for VPN nodes
- Cloudflare CDN for domain fronting
- Prometheus + Grafana for monitoring
- Loki for log aggregation

## Security

- All API traffic: TLS 1.3 minimum
- Internal services: mTLS (Istio)
- Keys stored in platform secure enclaves
- Certificate pinning on all clients
- Regular key rotation (7-90 days depending on type)
