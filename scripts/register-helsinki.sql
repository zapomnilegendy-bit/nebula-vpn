-- Register Helsinki node in the VPN servers table
-- Run when database is available: psql $DATABASE_URL -f register-helsinki.sql

INSERT INTO vpn_servers (name, host, ipv4, port, protocol, public_key, country, city, latitude, longitude, bandwidth_mbps, max_connections, tags)
VALUES (
    'Helsinki-FI-01',
    '83.147.216.161',
    '83.147.216.161',
    51820,
    'wireguard',
    'x4jFp1LCFM1PbixnflU0I6GJJjne5IaN5uyZLxOmdHE=',
    'FI',
    'Helsinki',
    60.1699,
    24.9384,
    1000,
    1000,
    ARRAY['FI', 'Helsinki', 'production']
)
ON CONFLICT DO NOTHING;
