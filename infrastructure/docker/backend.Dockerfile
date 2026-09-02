# Multi-stage build: security scan → build → distroless runtime
FROM golang:1.27-alpine AS builder

RUN apk add --no-cache git make gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE=api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /app/bin/${SERVICE} ./cmd/${SERVICE}

# Security scan stage
FROM aquasec/trivy:latest AS trivy
COPY --from=builder /app/bin /app/bin
RUN trivy filesystem --severity CRITICAL,HIGH --exit-code 1 --no-progress /app/bin || echo "⚠️  Trivy found issues (non-blocking)"

# Distroless runtime
FROM gcr.io/distroless/static:nonroot AS runtime

COPY --from=builder /app/bin/ /app/
COPY --from=builder /app/config/ /app/config/

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/app/api-gateway", "-healthcheck"]

EXPOSE 8080 9090

ENTRYPOINT ["/app/api-gateway"]
