use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VpnConfig {
    pub api_base_url: String,
    pub adaptive_transport: bool,
    pub auth_token: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ServerConfig {
    pub id: String,
    pub host: String,
    pub port: u16,
    pub public_key: String,
    pub protocol: String,
    pub country: String,
    pub city: String,
    pub load_percent: u8,
    pub latency_ms: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DeviceConfig {
    pub user_id: String,
    pub device_id: String,
    pub private_key: String,
    pub address: String,
    pub dns: Vec<String>,
    pub mtu: u16,
    pub split_tunnel: Option<SplitTunnelConfig>,
    pub kill_switch: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SplitTunnelConfig {
    pub enabled: bool,
    pub include_routes: Vec<String>,
    pub exclude_routes: Vec<String>,
}

impl VpnConfig {
    pub fn new(api_base_url: &str) -> Self {
        Self {
            api_base_url: api_base_url.to_string(),
            adaptive_transport: true,
            auth_token: None,
        }
    }

    pub fn with_auth_token(mut self, token: &str) -> Self {
        self.auth_token = Some(token.to_string());
        self
    }
}

impl ServerConfig {
    pub fn new(id: &str, host: &str, port: u16, public_key: &str) -> Self {
        Self {
            id: id.to_string(),
            host: host.to_string(),
            port,
            public_key: public_key.to_string(),
            protocol: "wireguard".to_string(),
            country: "DE".to_string(),
            city: "Frankfurt".to_string(),
            load_percent: 0,
            latency_ms: 0,
        }
    }
}
