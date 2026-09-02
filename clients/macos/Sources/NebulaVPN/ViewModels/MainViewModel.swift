import Foundation
import SwiftUI
import Combine
import LocalAuthentication

@MainActor
final class MainViewModel: ObservableObject {
    // MARK: - Published State
    @Published var connectionStatus: ConnectionStatus = .disconnected
    @Published var selectedServer: VPNServer?
    @Published var dataTransferred = DataUsage()
    @Published var connectionDuration: TimeInterval = 0
    @Published var activeProtocol: TransportProtocol = .wireGuard
    @Published var servers: [VPNServer] = VPNServer.samples
    @Published var settings: AppSettings
    @Published var isOnboarding = false
    @Published var showingAbout = false
    @Published var showingSettings = false
    @Published var showingServerPicker = false
    @Published var showingQRImporter = false
    @Published var isAuthenticated = false

    // MARK: - Private
    private let vpnService = VPNService()
    private let keychain = KeychainManager()
    private var durationTimer: Timer?
    private var trafficTimer: Timer?

    var isConnected: Bool {
        connectionStatus == .connected
    }

    // MARK: - Init
    init() {
        settings = keychain.loadSettings() ?? AppSettings()
        isOnboarding = !keychain.hasSeenOnboarding()
        authenticateUser()
    }

    // MARK: - Authentication
    func authenticateUser() {
        let context = LAContext()
        var error: NSError?

        guard context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error) else {
            isAuthenticated = true
            return
        }

        context.evaluatePolicy(.deviceOwnerAuthenticationWithBiometrics,
                                localizedReason: "Unlock NebulaVPN to access your secure connection") { success, _ in
            Task { @MainActor in
                self.isAuthenticated = success
                if !success { self.authenticateUser() }
            }
        }
    }

    // MARK: - Connection
    func connect() {
        guard let server = selectedServer else { return }
        connectionStatus = .connecting

        Task {
            do {
                try await vpnService.connect(to: server, protocol: activeProtocol)
                connectionStatus = .connected
                startMonitors()
            } catch {
                connectionStatus = .error(error.localizedDescription)
            }
        }
    }

    func disconnect() {
        stopMonitors()
        vpnService.disconnect()
        connectionStatus = .disconnected
        dataTransferred = DataUsage()
        connectionDuration = 0
    }

    func selectFastest() {
        selectedServer = servers
            .filter { $0.latencyMs > 0 }
            .min(by: { $0.latencyMs < $1.latencyMs })
    }

    func switchProtocol(to protocol: TransportProtocol) {
        activeProtocol = `protocol`
        if isConnected {
            Task { await reconnect() }
        }
    }

    private func reconnect() async {
        disconnect()
        try? await Task.sleep(nanoseconds: 500_000_000)
        connect()
    }

    // MARK: - Monitors
    private func startMonitors() {
        durationTimer = Timer.scheduledTimer(withTimeInterval: 1, repeats: true) { [weak self] _ in
            Task { @MainActor [weak self] in
                self?.connectionDuration += 1
            }
        }
        trafficTimer = Timer.scheduledTimer(withTimeInterval: 2, repeats: true) { [weak self] _ in
            Task { @MainActor [weak self] in
                guard let self = self, self.isConnected else { return }
                self.dataTransferred.download += Int64.random(in: 500...15000)
                self.dataTransferred.upload += Int64.random(in: 100...3000)
            }
        }
    }

    nonisolated func stopMonitors() {
        Task { @MainActor in
            self.durationTimer?.invalidate()
            self.trafficTimer?.invalidate()
            self.durationTimer = nil
            self.trafficTimer = nil
        }
    }

    var formattedDuration: String {
        let h = Int(connectionDuration) / 3600
        let m = (Int(connectionDuration) % 3600) / 60
        let s = Int(connectionDuration) % 60
        if h > 0 { return String(format: "%dh %dm", h, m) }
        return String(format: "%dm %ds", m, s)
    }

    // MARK: - QR Import
    func importFromQR(_ configString: String) {
        guard let server = parseServerConfig(configString) else { return }
        if !servers.contains(where: { $0.id == server.id }) {
            servers.append(server)
        }
        selectedServer = server
    }

    private func parseServerConfig(_ config: String) -> VPNServer? {
        return VPNServer.samples.first
    }

    // MARK: - Settings
    func saveSettings() {
        keychain.saveSettings(settings)
    }

    func completeOnboarding() {
        keychain.markOnboardingSeen()
        isOnboarding = false
    }

    deinit {
        stopMonitors()
    }
}
