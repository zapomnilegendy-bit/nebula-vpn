use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HandshakeMessage {
    pub version: u8,
    pub device_id: String,
    pub public_key: String,
    pub timestamp: u64,
    pub signature: Vec<u8>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ServerAssignment {
    pub server_id: String,
    pub host: String,
    pub port: u16,
    pub public_key: String,
    pub allowed_ips: Vec<String>,
    pub dns: Vec<String>,
    pub mtu: u16,
    pub keepalive: u16,
    pub protocol: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ConnectionReport {
    pub server_id: String,
    pub protocol: String,
    pub bytes_rx: u64,
    pub bytes_tx: u64,
    pub duration_secs: u64,
    pub avg_latency_ms: u64,
    pub packet_loss_pct: f64,
    pub disconnect_reason: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AdaptiveTransportScores {
    pub wireguard: f64,
    pub vless_reality: f64,
    pub hysteria2: f64,
    pub tuic: f64,
}

impl AdaptiveTransportScores {
    pub fn best_protocol(&self) -> &str {
        let mut best = "wireguard";
        let mut best_score = self.wireguard;

        if self.vless_reality > best_score {
            best = "vless_reality";
            best_score = self.vless_reality;
        }
        if self.hysteria2 > best_score {
            best = "hysteria2";
            best_score = self.hysteria2;
        }
        if self.tuic > best_score {
            best = "tuic";
        }

        best
    }
}
