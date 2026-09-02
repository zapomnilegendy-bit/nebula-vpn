#!/usr/bin/env bash
set -euo pipefail

echo "🔒 NebulaVPN Security Audit"
echo "=========================="

# Go security
echo ""
echo "→ Go security check (gosec)..."
cd backend
gosec -no-fail -fmt=text ./...
cd -

echo ""
echo "→ Go vulnerability check (govulncheck)..."
cd backend
govulncheck ./...
cd -

# Docker security
echo ""
echo "→ Dockerfile lint (hadolint)..."
hadolint infrastructure/docker/backend.Dockerfile 2>/dev/null || echo "  (hadolint not installed)"

# Dependency audit
echo ""
echo "→ NPM audit..."
cd web/user-panel && npm audit --audit-level=high 2>/dev/null || true; cd -
cd web/admin-panel && npm audit --audit-level=high 2>/dev/null || true; cd -

# Secrets scan
echo ""
echo "→ Scanning for secrets in repo..."
git grep -l -i 'password\s*=\|secret\s*=\|api_key\s*=\|token\s*=' -- '*.go' '*.ts' '*.tsx' '*.swift' '*.kt' '*.rs' 2>/dev/null || true

echo ""
echo "✅ Audit complete"
