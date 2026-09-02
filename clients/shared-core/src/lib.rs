pub mod transport;
pub mod config;
pub mod crypto;
pub mod protocol;

use std::sync::Arc;
use tokio::sync::RwLock;
use transport::ConnectionState;

/// Core VPN client engine
pub struct NebulaCore {
    state: Arc<RwLock<EngineState>>,
    config: config::VpnConfig,
}

#[derive(Debug, Clone)]
pub struct EngineState {
    pub connection_state: ConnectionState,
    pub active_protocol: transport::Protocol,
    pub connected_server: Option<String>,
    pub bytes_rx: u64,
    pub bytes_tx: u64,
    pub latency_ms: Option<u64>,
}

impl NebulaCore {
    pub fn new(config: config::VpnConfig) -> Self {
        Self {
            state: Arc::new(RwLock::new(EngineState {
                connection_state: ConnectionState::Disconnected,
                active_protocol: transport::Protocol::WireGuard,
                connected_server: None,
                bytes_rx: 0,
                bytes_tx: 0,
                latency_ms: None,
            })),
            config,
        }
    }

    pub async fn connect(&self, server: &config::ServerConfig) -> Result<(), Box<dyn std::error::Error>> {
        let mut state = self.state.write().await;
        state.connection_state = ConnectionState::Connecting;
        state.connected_server = Some(server.host.clone());

        // Protocol selection
        let transport = self.select_transport(server).await?;
        transport.connect().await?;

        state.connection_state = ConnectionState::Connected;
        state.active_protocol = transport.protocol();
        Ok(())
    }

    pub async fn disconnect(&self) -> Result<(), Box<dyn std::error::Error>> {
        let mut state = self.state.write().await;
        state.connection_state = ConnectionState::Disconnected;
        state.connected_server = None;
        Ok(())
    }

    pub async fn get_state(&self) -> EngineState {
        self.state.read().await.clone()
    }

    async fn select_transport(&self, server: &config::ServerConfig) -> Result<Box<dyn transport::Transport>, Box<dyn std::error::Error>> {
        // Adaptive transport selection based on server capabilities and network conditions
        if self.config.adaptive_transport {
            // TODO: score each protocol based on latency/packet loss
        }

        Ok(Box::new(transport::WireGuardTransport::new(
            server.host.clone(),
            server.port,
            server.public_key.clone(),
        )?))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_engine_creation() {
        let config = config::VpnConfig {
            api_base_url: "https://api.nebula-vpn.io".to_string(),
            adaptive_transport: true,
            auth_token: None,
        };
        let engine = NebulaCore::new(config);
        assert!(engine.state.try_read().is_ok());
    }
}
