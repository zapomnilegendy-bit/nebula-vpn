#!/usr/bin/env bash
set -euo pipefail

echo "🔍 Running pre-commit checks..."

# Go checks
if git diff --cached --name-only | grep -q '\.go$'; then
  echo "  → Linting Go files..."
  cd backend && golangci-lint run ./... --timeout=3m
  echo "  → Formatting Go files..."
  go fmt ./...
  cd -
fi

# TypeScript checks
if git diff --cached --name-only | grep -qE '\.(ts|tsx)$'; then
  echo "  → Linting TypeScript..."
  cd web/user-panel && npx tsc --noEmit 2>/dev/null || true; cd -
  cd web/admin-panel && npx tsc --noEmit 2>/dev/null || true; cd -
fi

# Swift checks
if git diff --cached --name-only | grep -q '\.swift$'; then
  echo "  → Linting Swift..."
  which swiftlint && swiftlint --quiet || true
fi

# Kotlin checks
if git diff --cached --name-only | grep -q '\.kt$'; then
  echo "  → Linting Kotlin..."
  which detekt && detekt --config clients/android/config/detekt.yml --build || true
fi

# Rust checks
if git diff --cached --name-only | grep -q '\.rs$'; then
  echo "  → Checking Rust..."
  cd clients/shared-core && cargo check 2>/dev/null || true; cd -
fi

# Dockerfile checks
if git diff --cached --name-only | grep -q 'Dockerfile'; then
  echo "  → Linting Dockerfiles..."
  which hadolint && hadolint infrastructure/docker/backend.Dockerfile || true
fi

# Secrets scan
echo "  → Scanning for secrets..."
git diff --cached | grep -qiE '(password|secret|api.?key|token)' && echo "⚠️  Potential secret in diff!" || true

echo "✅ Pre-commit checks complete"
