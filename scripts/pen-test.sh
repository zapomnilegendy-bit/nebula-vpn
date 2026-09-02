#!/usr/bin/env bash
# Nebula VPN — Penetration Test Suite
# Tests: auth brute force, JWT fuzzing, SQL injection, rate limiting, CORS
# Usage: ./pen-test.sh [target_url] [--verbose]

set -euo pipefail

TARGET="${1:-http://localhost:8080}"
VERBOSE=false
[[ "${2:-}" == "--verbose" ]] && VERBOSE=true

PASS=0
FAIL=0

log_pass()  { echo "  ✅ $1"; ((PASS++)); }
log_fail()  { echo "  ❌ $1"; ((FAIL++)); }
log_info()  { echo "  ℹ️  $1"; }
banner()    { echo -e "\n═══ $1 ═══"; }

cleanup() { rm -f /tmp/pen-test-*.txt; }
trap cleanup EXIT

# ── Auth Brute Force ──────────────────────────────────────────────
banner "Auth: Brute Force Protection"
for i in $(seq 1 10); do
    code=$(curl -s -o /dev/null -w "%{http_code}" \
        -X POST "$TARGET/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"test@test.com","password":"wrongpass123"}')
    if [ "$code" = "429" ]; then
        log_pass "Rate limited after $i attempts (HTTP 429)"
        break
    fi
done
if [ "$code" != "429" ]; then
    log_fail "No rate limiting detected after 10 attempts (last: $code)"
fi

# ── SQL Injection ──────────────────────────────────────────────────
banner "Security: SQL Injection"
for payload in \
    "' OR '1'='1" \
    "'; DROP TABLE users; --" \
    "' UNION SELECT * FROM users--" \
    "admin'--" \
    "1; SELECT pg_sleep(5)--"; do

    encoded=$(python3 -c "import urllib.parse; print(urllib.parse.quote('$payload'))")
    code=$(curl -s -o /dev/null -w "%{http_code}" \
        "$TARGET/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$payload@test.com\",\"password\":\"test123\"}")
    if [ "$code" = "500" ]; then
        log_fail "SQL injection may have succeeded: $payload"
        continue 2
    fi
done
log_pass "SQL injection inputs rejected"

# ── JWT Fuzzing ────────────────────────────────────────────────────
banner "Security: JWT Fuzzing"
for token in \
    "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjM0NTY3ODkwIn0." \
    "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.WRONG" \
    "invalid-token" \
    "$(echo '{}' | base64).$(echo '{"sub":"1","role":"admin"}' | base64).dGVzdA"; do

    code=$(curl -s -o /dev/null -w "%{http_code}" \
        "$TARGET/v1/devices" \
        -H "Authorization: Bearer $token")
    if [ "$code" = "200" ]; then
        log_fail "JWT token accepted: $token"
        continue 2
    fi
done
log_pass "JWT fuzzing tokens rejected"

# ── CORS ───────────────────────────────────────────────────────────
banner "Security: CORS"
code=$(curl -s -o /dev/null -w "%{http_code}" \
    -X OPTIONS "$TARGET/v1/health" \
    -H "Origin: https://evil.com" \
    -H "Access-Control-Request-Method: GET" 2>/dev/null || true)
origin=$(curl -s -D- "$TARGET/v1/health" \
    -H "Origin: https://evil.com" 2>/dev/null \
    | grep -i "access-control-allow-origin" | tr -d '\r' || true)
if echo "$origin" | grep -qi "evil.com"; then
    log_fail "CORS allows evil.com origin"
else
    log_pass "CORS blocked evil.com origin"
fi

# ── Input Validation ──────────────────────────────────────────────
banner "Security: Input Validation"
# XSS in register
code=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$TARGET/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"email":"<script>alert(1)</script>@test.com","password":"Pass1234"}')
if [ "$code" = "400" ]; then
    log_pass "XSS email rejected"
else
    log_fail "XSS email accepted (HTTP $code)"
fi

# Oversized payload
payload=$(python3 -c "print('a'*1000000)")
code=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$TARGET/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"test@test.com\",\"password\":\"$payload\"}")
if [ "$code" = "413" ] || [ "$code" = "400" ]; then
    log_pass "Oversized payload rejected (HTTP $code)"
else
    log_fail "Oversized payload accepted (HTTP $code)"
fi

# ── Security Headers ──────────────────────────────────────────────
banner "Security: Response Headers"
headers=$(curl -s -I "$TARGET/v1/health" 2>/dev/null)
for header in \
    "X-Content-Type-Options: nosniff" \
    "X-Frame-Options: DENY" \
    "X-XSS-Protection: 1; mode=block" \
    "Referrer-Policy" \
    "Permissions-Policy"; do
    if echo "$headers" | grep -qi "$header"; then
        log_pass "Header present: $header"
    else
        log_fail "Header missing: $header"
    fi
done

# ── Summary ────────────────────────────────────────────────────────
banner "Results"
echo "  Passed: $PASS"
echo "  Failed: $FAIL"
[ "$FAIL" -eq 0 ] && echo "  Result: ✅ All tests passed" || echo "  Result: ❌ $FAIL tests failed"
