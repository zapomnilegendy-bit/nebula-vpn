#!/usr/bin/env bash
set -euo pipefail

echo "🔧 Installing development tools..."

# Go tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

# Node tools
npm install -g typescript prettier eslint

# Swift tools
if command -v brew &>/dev/null; then
  brew install swiftlint
fi

# Docker tools
if command -v brew &>/dev/null; then
  brew install hadolint
fi

# Pre-commit hook
cp scripts/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
echo "  → Installed pre-commit hook"

echo "✅ All tools installed"
echo ""
echo "Run individually:"
echo "  make lint    — Go lint"
echo "  make test    — Go tests"
echo "  make audit   — Security audit"
