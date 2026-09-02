import NetworkExtension
import os

class PacketTunnelProvider: NEPacketTunnelProvider {
    private let log = Logger(subsystem: "com.nebula.vpn.macos", category: "PacketTunnel")
    private var cfg: [String: Any]?

    override func startTunnel(options: [String: NSObject]? = nil) async throws {
        log.info("Starting tunnel with options: \(options?.description ?? "none")")

        cfg = (protocolConfiguration as? NETunnelProviderProtocol)?.providerConfiguration

        // Read remote address from config or use default
        let remoteAddress = (cfg?["host"] as? String) ?? "10.0.0.1"
        let mtu = (cfg?["mtu"] as? Int) ?? 1420

        let settings = NEPacketTunnelNetworkSettings(tunnelRemoteAddress: remoteAddress)

        // MTU
        settings.mtu = NSNumber(integerLiteral: mtu)

        // IPv4
        let ipv4 = NEIPv4Settings(addresses: ["10.0.0.2"], subnetMasks: ["255.255.255.0"])
        ipv4.includedRoutes = [NEIPv4Route.default()]
        ipv4.excludedRoutes = [
            NEIPv4Route(destinationAddress: "10.0.0.0", subnetMask: "255.0.0.0"),
            NEIPv4Route(destinationAddress: "172.16.0.0", subnetMask: "255.240.0.0"),
            NEIPv4Route(destinationAddress: "192.168.0.0", subnetMask: "255.255.0.0"),
            NEIPv4Route(destinationAddress: "100.64.0.0", subnetMask: "255.192.0.0"),
        ]
        settings.ipv4Settings = ipv4

        // IPv6
        let ipv6 = NEIPv6Settings(addresses: ["fd00::2"], networkPrefixLengths: [64])
        ipv6.includedRoutes = [NEIPv6Route.default()]
        settings.ipv6Settings = ipv6

        // DNS
        let dns = NEDNSSettings(servers: ["1.1.1.1", "1.0.0.1"])
        dns.matchDomains = [""] // Route all DNS through tunnel
        settings.dnsSettings = dns

        // Apply settings
        try await setTunnelNetworkSettings(settings)
        log.info("Tunnel settings applied successfully")

        // Start packet processing
        startPacketFlow()
    }

    override func stopTunnel(with reason: NEProviderStopReason) async {
        log.info("Stopping tunnel. Reason: \(reason.rawValue)")
    }

    override func handleAppMessage(_ messageData: Data) async -> Data? {
        log.info("Received app message: \(messageData.count) bytes")
        switch String(data: messageData, encoding: .utf8) {
        case "get-status":
            return Data("connected".utf8)
        case "get-stats":
            let stats = """
            {"rx": 0, "tx": 0}
            """
            return Data(stats.utf8)
        default:
            return nil
        }
    }

    // MARK: - Packet Flow
    private func startPacketFlow() {
        Task { [weak self] in
            guard let self = self else { return }
            let flow = self.packetFlow
            var isRunning = true

            while isRunning {
                let packets = await flow.readPackets()
                for (data, protocolFamily) in packets {
                    // In production, encrypt/decrypt via Rust core
                    // For now, pass through with obfuscation padding
                    let processed = data + Data(repeating: 0, count: Int.random(in: 40...200))
                    flow.writePackets([processed], withProtocols: [protocolFamily])
                }

                // Small sleep to prevent tight loop
                try? await Task.sleep(nanoseconds: 10_000_000)
            }
        }
    }
}
