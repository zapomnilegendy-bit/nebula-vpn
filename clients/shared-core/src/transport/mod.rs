use async_trait::async_trait;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Protocol {
    WireGuard,
    VLESSReality,
    Hysteria2,
    TUIC,
}

impl Protocol {
    pub fn as_str(&self) -> &'static str {
        match self {
            Protocol::WireGuard => "wireguard",
            Protocol::VLESSReality => "vless_reality",
            Protocol::Hysteria2 => "hysteria2",
            Protocol::TUIC => "tuic",
        }
    }

    pub fn from_str(s: &str) -> Option<Self> {
        match s {
            "wireguard" => Some(Protocol::WireGuard),
            "vless_reality" => Some(Protocol::VLESSReality),
            "hysteria2" => Some(Protocol::Hysteria2),
            "tuic" => Some(Protocol::TUIC),
            _ => None,
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ConnectionState {
    Disconnected,
    Connecting,
    Connected,
    Reconnecting,
    Error,
}

#[async_trait]
pub trait Transport: Send + Sync {
    fn protocol(&self) -> Protocol;
    async fn connect(&self) -> Result<(), Box<dyn std::error::Error>>;
    async fn disconnect(&self) -> Result<(), Box<dyn std::error::Error>>;
    async fn latency(&self) -> Option<u64>;
    fn is_healthy(&self) -> bool;
}

pub struct WireGuardTransport {
    endpoint: String,
    port: u16,
    public_key: String,
    private_key: Option<String>,
}

impl WireGuardTransport {
    pub fn new(endpoint: String, port: u16, public_key: String) -> Result<Self, Box<dyn std::error::Error>> {
        Ok(Self {
            endpoint,
            port,
            public_key,
            private_key: None,
        })
    }
}

#[async_trait]
impl Transport for WireGuardTransport {
    fn protocol(&self) -> Protocol {
        Protocol::WireGuard
    }

    async fn connect(&self) -> Result<(), Box<dyn std::error::Error>> {
        tracing::info!("connecting WireGuard to {}:{}", self.endpoint, self.port);
        Ok(())
    }

    async fn disconnect(&self) -> Result<(), Box<dyn std::error::Error>> {
        tracing::info!("disconnecting WireGuard");
        Ok(())
    }

    async fn latency(&self) -> Option<u64> {
        None
    }

    fn is_healthy(&self) -> bool {
        true
    }
}
