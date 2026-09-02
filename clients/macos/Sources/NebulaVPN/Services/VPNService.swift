import Foundation
import NetworkExtension

final class VPNService {
    private var manager: NETunnelProviderManager?
    private let bundleID = "com.nebula.vpn.macos.PacketTunnel"

    // MARK: - Setup
    private func setupManager() async throws {
        let managers = try await NETunnelProviderManager.loadAllFromPreferences()

        if let existing = managers.first {
            manager = existing
        } else {
            let newManager = NETunnelProviderManager()
            let proto = NETunnelProviderProtocol()
            proto.providerBundleIdentifier = bundleID
            proto.serverAddress = "nebula-vpn.io"
            newManager.protocolConfiguration = proto
            newManager.localizedDescription = "NebulaVPN"
            manager = newManager
        }

        try await manager?.saveToPreferences()
    }

    // MARK: - Connect
    func connect(to server: VPNServer, protocol: TransportProtocol) async throws {
        try await setupManager()

        guard let manager = manager else {
            throw VPNError.notConfigured
        }

        guard let proto = manager.protocolConfiguration as? NETunnelProviderProtocol else {
            throw VPNError.notConfigured
        }

        let config: [String: Any] = [
            "server_id": server.id,
            "host": server.host,
            "port": server.port,
            "public_key": server.publicKey,
            "protocol": `protocol`.rawValue,
            "mtu": 1420,
        ]

        proto.providerConfiguration = config
        manager.protocolConfiguration = proto
        try await manager.saveToPreferences()
        try await manager.loadFromPreferences()

        do {
            try manager.connection.startVPNTunnel(options: [:])
        } catch {
            throw VPNError.startFailed(error)
        }
    }

    // MARK: - Disconnect
    func disconnect() {
        manager?.connection.stopVPNTunnel()
    }

    // MARK: - Status
    var status: NEVPNStatus {
        manager?.connection.status ?? .invalid
    }
}

// MARK: - Errors
enum VPNError: Swift.Error, LocalizedError {
    case notConfigured
    case startFailed(Swift.Error)
    case missingEntitlement

    var errorDescription: String? {
        switch self {
        case .notConfigured:
            return "VPN not configured. Please try again."
        case .startFailed(let error):
            return "Failed to start: \(error.localizedDescription)"
        case .missingEntitlement:
            return "Missing Network Extension entitlement."
        }
    }
}
