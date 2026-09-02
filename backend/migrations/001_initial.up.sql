CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- Devices
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    platform VARCHAR(20) NOT NULL CHECK (platform IN ('macos', 'windows', 'ios', 'android')),
    public_key TEXT NOT NULL,
    private_key_encrypted TEXT NOT NULL,
    ip_address INET NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_connected_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_devices_user ON devices(user_id);
CREATE INDEX idx_devices_platform ON devices(platform);

-- VPN Servers
CREATE TABLE IF NOT EXISTS vpn_servers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    ipv4 INET,
    ipv6 INET,
    port INTEGER NOT NULL,
    protocol VARCHAR(20) NOT NULL,
    public_key TEXT NOT NULL,
    country VARCHAR(2) NOT NULL,
    city VARCHAR(100),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    load_percent INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_healthy BOOLEAN DEFAULT true,
    bandwidth_mbps INTEGER DEFAULT 1000,
    current_connections INTEGER DEFAULT 0,
    max_connections INTEGER DEFAULT 1000,
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_servers_country ON vpn_servers(country);
CREATE INDEX idx_servers_protocol ON vpn_servers(protocol);
CREATE INDEX idx_servers_healthy ON vpn_servers(is_healthy) WHERE is_healthy = true;

-- Plans
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price_cents INTEGER NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    max_devices INTEGER NOT NULL,
    max_speed_mbps INTEGER,
    protocols_enabled TEXT[] DEFAULT '{"wireguard"}',
    locations_enabled TEXT[] DEFAULT '{}',
    trial_period_days INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES plans(id),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'canceled', 'expired')),
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    auto_renew BOOLEAN DEFAULT true,
    stripe_subscription_id VARCHAR(255),
    apple_receipt TEXT,
    google_purchase_token TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_expires ON subscriptions(expires_at) WHERE status = 'active';

-- Connections (session log)
CREATE TABLE IF NOT EXISTS connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    device_id UUID REFERENCES devices(id),
    server_id UUID REFERENCES vpn_servers(id),
    protocol VARCHAR(20) NOT NULL,
    bytes_rx BIGINT DEFAULT 0,
    bytes_tx BIGINT DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    reason VARCHAR(50)
);

CREATE INDEX idx_connections_user ON connections(user_id);
CREATE INDEX idx_connections_time ON connections(started_at DESC);

-- Promocodes
CREATE TABLE IF NOT EXISTS promocodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    discount_percent INTEGER NOT NULL CHECK (discount_percent > 0 AND discount_percent <= 100),
    max_uses INTEGER,
    current_uses INTEGER DEFAULT 0,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_promocodes_code ON promocodes(code);

-- Insert default plans
INSERT INTO plans (name, description, price_cents, currency, max_devices, max_speed_mbps, protocols_enabled, locations_enabled, trial_period_days) VALUES
    ('Free', 'Basic VPN access with limited features', 0, 'USD', 1, 50, '{"wireguard"}', '{}', 0),
    ('Pro', 'Full VPN access with all protocols', 999, 'USD', 5, 500, '{"wireguard","vless_reality","hysteria2"}', '{}', 7),
    ('Ultimate', 'Unlimited speed and all locations', 1999, 'USD', 10, NULL, '{"wireguard","vless_reality","hysteria2","tuic"}', '{}', 14)
ON CONFLICT DO NOTHING;
