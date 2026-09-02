import Foundation
import SwiftUI

// MARK: - Connection State
enum ConnectionStatus: Equatable {
    case disconnected
    case connecting
    case connected
    case reconnecting
    case error(String)

    var label: String {
        switch self {
        case .disconnected: return "Not Connected"
        case .connecting:   return "Connecting…"
        case .connected:    return "Connected"
        case .reconnecting: return "Reconnecting…"
        case .error(let m): return m
        }
    }

    var color: Color {
        switch self {
        case .disconnected: return .secondary.opacity(0.35)
        case .connecting:   return .orange
        case .connected:    return Color(red: 0.22, green: 0.78, blue: 0.35)
        case .reconnecting: return .orange
        case .error:        return .red
        }
    }

    var icon: String {
        switch self {
        case .disconnected: return "shield.slash"
        case .connecting:   return "shield"
        case .connected:    return "shield.checkered"
        case .reconnecting: return "arrow.triangle.2.circlepath"
        case .error:        return "exclamationmark.triangle"
        }
    }
}

// MARK: - Transport Protocol
enum TransportProtocol: String, CaseIterable, Identifiable, Codable {
    case wireGuard      = "wireguard"
    case vlessReality   = "vless_reality"
    case hysteria2      = "hysteria2"
    case tuic           = "tuic"

    var id: String { rawValue }

    var displayName: String {
        switch self {
        case .wireGuard:    return "WireGuard"
        case .vlessReality: return "VLESS+Reality"
        case .hysteria2:    return "Hysteria2"
        case .tuic:         return "TUIC"
        }
    }

    var priority: Int {
        switch self {
        case .wireGuard:    return 0
        case .vlessReality: return 1
        case .hysteria2:    return 2
        case .tuic:         return 3
        }
    }
}

// MARK: - Server Model
struct VPNServer: Identifiable, Codable, Equatable {
    let id: String
    let name: String
    let city: String
    let country: String
    let host: String
    let port: Int
    let publicKey: String
    let loadPercent: Int
    let latencyMs: Int
    let isPremium: Bool
    let supportedProtocols: [TransportProtocol]

    var displayName: String { "\(city), \(countryFlag)" }

    var countryFlag: String {
        let base: UInt32 = 127462
        let scalars = country.uppercased().unicodeScalars
        guard scalars.count == 2 else { return "🏳️" }
        return String(String.UnicodeScalarView(
            scalars.compactMap { UnicodeScalar(base + $0.value - 65) }
        ))
    }

    var loadColor: Color {
        if loadPercent < 40 { return .green }
        if loadPercent < 70 { return .orange }
        return .red
    }

    var detailString: String {
        "\(latencyMs)ms · \(loadPercent)% · \(supportedProtocols.first?.displayName ?? "WireGuard")"
    }
}

// MARK: - Data Usage
struct DataUsage: Equatable {
    var download: Int64 = 0
    var upload: Int64 = 0

    var formattedDownload: String {
        ByteCountFormatter.string(fromByteCount: download, countStyle: .file)
    }
    var formattedUpload: String {
        ByteCountFormatter.string(fromByteCount: upload, countStyle: .file)
    }
}

// MARK: - App Settings
struct AppSettings: Codable {
    var launchAtLogin: Bool = true
    var killSwitch: Bool = true
    var autoConnect: Bool = false
    var adaptiveProtocol: Bool = true
    var dnsOverHTTPS: Bool = true
    var dnsOverTLS: Bool = true
    var selectedProtocol: TransportProtocol = .wireGuard
    var appearance: Appearance = .system
    var selectedServerID: String?
}

enum Appearance: String, CaseIterable, Codable {
    case system, light, dark
}

// MARK: - Sample Data
extension VPNServer {
    static let samples: [VPNServer] = [
        .init(id: "eu-fra-01", name: "fra1.nebula", city: "Frankfurt", country: "DE",
              host: "fra1.nebula-vpn.io", port: 51820, publicKey: "key1", loadPercent: 42, latencyMs: 12,
              isPremium: false, supportedProtocols: TransportProtocol.allCases),
        .init(id: "eu-ams-01", name: "ams1.nebula", city: "Amsterdam", country: "NL",
              host: "ams1.nebula-vpn.io", port: 51820, publicKey: "key2", loadPercent: 58, latencyMs: 18,
              isPremium: false, supportedProtocols: [.wireGuard, .vlessReality, .hysteria2]),
        .init(id: "gb-lon-01", name: "lon1.nebula", city: "London", country: "GB",
              host: "lon1.nebula-vpn.io", port: 51820, publicKey: "key6", loadPercent: 55, latencyMs: 24,
              isPremium: false, supportedProtocols: [.wireGuard, .vlessReality]),
        .init(id: "us-nyc-01", name: "nyc1.nebula", city: "New York", country: "US",
              host: "nyc1.nebula-vpn.io", port: 51820, publicKey: "key3", loadPercent: 71, latencyMs: 89,
              isPremium: true, supportedProtocols: TransportProtocol.allCases),
        .init(id: "us-sfo-01", name: "sfo1.nebula", city: "San Francisco", country: "US",
              host: "sfo1.nebula-vpn.io", port: 51820, publicKey: "key4", loadPercent: 34, latencyMs: 105,
              isPremium: false, supportedProtocols: [.wireGuard, .hysteria2]),
        .init(id: "jp-tokyo-01", name: "tky1.nebula", city: "Tokyo", country: "JP",
              host: "tky1.nebula-vpn.io", port: 51820, publicKey: "key7", loadPercent: 23, latencyMs: 185,
              isPremium: true, supportedProtocols: [.wireGuard, .vlessReality, .tuic]),
    ]
}
