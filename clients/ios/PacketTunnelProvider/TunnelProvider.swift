import NetworkExtension
import os.log

/// Battery-optimized packet tunnel provider with batched processing and idle sensing
final class TunnelProvider: NEPacketTunnelProvider {

    private var idleSince: Date?
    private var packetBatch: [Data] = []
    private let batchSize = 32
    private let maxBatchDelay: TimeInterval = 0.01 // 10ms
    private var lastFlush = Date()
    private let logger = OSLog(subsystem: "com.nebula.vpn", category: "TunnelProvider")

    override func startTunnel(options: [String : NSObject]?, completionHandler: @escaping (Error?) -> Void) {
        let settings = NEPacketTunnelNetworkSettings(tunnelRemoteAddress: "10.0.0.1")

        // IPv4
        let ipv4 = NEIPv4Settings(addresses: ["10.0.0.2"], subnetMasks: ["255.255.255.0"])
        ipv4.includedRoutes = [NEIPv4Route.default()]
        settings.ipv4Settings = ipv4

        // IPv6
        let ipv6 = NEIPv6Settings(addresses: ["fd00::2"], networkPrefixLengths: [64])
        ipv6.includedRoutes = [NEIPv6Route.default()]
        settings.ipv6Settings = ipv6

        // DNS
        let dns = NEDNSSettings(servers: ["1.1.1.1", "8.8.8.8"])
        dns.matchDomains = [""]
        settings.dnssettings = dns

        // MTU
        settings.mtu = 1400

        setTunnelNetworkSettings(settings) { [weak self] error in
            if let error = error {
                os_log(.error, log: self?.logger ?? .default, "Failed to set tunnel settings: %@", error.localizedDescription)
                completionHandler(error)
                return
            }
            self?.startPacketHandlers()
            os_log(.info, log: self?.logger ?? .default, "Tunnel started successfully")
            completionHandler(nil)
        }
    }

    override func stopTunnel(with reason: NEProviderStopReason, completionHandler: @escaping () -> Void) {
        flushPacketBatch()
        os_log(.info, log: logger, "Tunnel stopped: %d", reason.rawValue)
        completionHandler()
    }

    private func startPacketHandlers() {
        // Read packets with idle-aware batching
        packetFlow.readPackets { [weak self] packets, protocols in
            guard let self = self else { return }

            for (index, packet) in packets.enumerated() {
                let proto = protocols[index].intValue
                self.enqueuePacket(packet, protocol: proto)
            }

            // Continue reading
            self.startPacketHandlers()
        }
    }

    // MARK: - Batched Packet Processing

    private func enqueuePacket(_ packet: Data, protocol proto: Int) {
        packetBatch.append(packet)

        if packetBatch.count >= batchSize || Date().timeIntervalSince(lastFlush) >= maxBatchDelay {
            flushPacketBatch()
        }

        // Track idle state
        idleSince = nil
    }

    private func flushPacketBatch() {
        guard !packetBatch.isEmpty else { return }
        lastFlush = Date()

        // Process batched packets
        var writePackets: [Data] = []
        for packet in packetBatch {
            if let encrypted = processPacket(packet) {
                writePackets.append(encrypted)
            }
        }
        packetBatch.removeAll(keepingCapacity: true)

        if !writePackets.isEmpty {
            let protocols = writePackets.map { NSNumber(value: AF_INET) }
            packetFlow.writePackets(writePackets, withProtocols: protocols)
        }
    }

    private func processPacket(_ packet: Data) -> Data? {
        // In production: encrypt/decrypt via Rust core
        // For now, pass through with obfuscation padding
        var processed = packet
        processed.append(contentsOf: generatePadding())
        return processed
    }

    private func generatePadding() -> [UInt8] {
        // Random padding 40-200 bytes to defeat traffic analysis
        let length = Int.random(in: 40...200)
        return (0..<length).map { _ in UInt8.random(in: 0...255) }
    }

    override func handleAppMessage(_ messageData: Data, completionHandler: ((Data?) -> Void)?) {
        // Handle messages from the app
        completionHandler?(nil)
    }

    override func sleep(completionHandler: @escaping () -> Void) {
        flushPacketBatch()
        idleSince = Date()
        completionHandler()
    }

    override func wake() {
        idleSince = nil
        startPacketHandlers()
    }
}
