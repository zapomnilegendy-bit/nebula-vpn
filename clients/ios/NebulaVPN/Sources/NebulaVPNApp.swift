import SwiftUI
import UserNotifications

@main
struct NebulaVPNApp: App {
    @UIApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    @StateObject private var appState = AppState()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(appState)
                .toast(toast: $appState.toast)
                .confirmAlert(alert: $appState.confirmAlert)
                .onOpenURL { url in
                    handleDeepLink(url)
                }
                .onAppear {
                    appState.checkOnboarding()
                }
        }
    }

    private func handleDeepLink(_ url: URL) {
        guard let route = DeepLinkHandler.shared.handle(url) else { return }
        appState.handleDeepLink(route)
    }
}

// MARK: - App Delegate
final class AppDelegate: NSObject, UIApplicationDelegate {
    let backgroundManager = BackgroundManager.shared
    let batteryManager = BatteryManager.shared
    let crashReporter = CrashReporter.shared

    func application(
        _ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]? = nil
    ) -> Bool {
        crashReporter.start()
        batteryManager.start()
        Task { @MainActor in
            backgroundManager.registerTasks()
            backgroundManager.scheduleTasks()
            await backgroundManager.registerForPushNotifications()
        }
        return true
    }

    func application(
        _ app: UIApplication,
        open url: URL,
        options: [UIApplication.OpenURLOptionsKey: Any] = [:]
    ) -> Bool {
        return DeepLinkHandler.shared.handle(url) != nil
    }

    func application(
        _ application: UIApplication,
        continue userActivity: NSUserActivity,
        restorationHandler: @escaping ([UIUserActivityRestoring]) -> Void
    ) -> Bool {
        guard userActivity.activityType == NSUserActivityTypeBrowsingWeb,
              let url = userActivity.webpageURL else { return false }
        return DeepLinkHandler.shared.handle(url) != nil
    }
}

// MARK: - AppState
@MainActor
class AppState: ObservableObject {
    @Published var isConnected = false
    @Published var selectedServer: VPNServer?
    @Published var connectionStatus: ConnectionStatus = .disconnected
    @Published var showingServerList = false
    @Published var dataTransferred = DataUsage()
    @Published var connectionTime: TimeInterval = 0
    @Published var showingSettings = false
    @Published var showingAbout = false
    @Published var isOnboarding = false
    @Published var toast: Toast?
    @Published var confirmAlert: ConfirmAlert?

    @AppStorage("hasSeenOnboarding") private var hasSeenOnboarding = false

    func checkOnboarding() {
        isOnboarding = !hasSeenOnboarding
    }

    func completeOnboarding() {
        hasSeenOnboarding = true
        isOnboarding = false
    }

    var statusColor: Color {
        switch connectionStatus {
        case .disconnected: return .secondary.opacity(0.3)
        case .connecting: return .orange
        case .connected: return Color(red: 0.22, green: 0.78, blue: 0.35)
        case .reconnecting: return .orange
        case .error: return .red
        }
    }

    var statusIcon: String {
        switch connectionStatus {
        case .disconnected: return "shield"
        case .connecting: return "shield"
        case .connected: return "shield.checkered"
        case .reconnecting: return "shield"
        case .error: return "exclamationmark.shield"
        }
    }

    var statusText: String {
        switch connectionStatus {
        case .disconnected: return "Not Connected"
        case .connecting: return "Connecting…"
        case .connected: return "Connected"
        case .reconnecting: return "Reconnecting…"
        case .error(let msg): return msg
        }
    }

    func connect() {
        guard selectedServer != nil else {
            showToast("Select a server first", type: .warning)
            return
        }
        connectionStatus = .connecting
        Task {
            try? await Task.sleep(nanoseconds: 600_000_000)
            withAnimation(.spring(response: 0.6)) {
                isConnected = true
                connectionStatus = .connected
            }
            showToast("Connected to \(selectedServer?.city ?? "VPN")", type: .success)
            startConnectionTimer()
            simulateTraffic()
        }
    }

    func disconnect() {
        withAnimation(.spring(response: 0.35)) {
            isConnected = false
            connectionStatus = .disconnected
            dataTransferred = DataUsage()
            connectionTime = 0
        }
        showToast("Disconnected", type: .info)
    }

    private func startConnectionTimer() {
        Task {
            while isConnected {
                try? await Task.sleep(nanoseconds: 1_000_000_000)
                connectionTime += 1
            }
        }
    }

    private func simulateTraffic() {
        Task {
            while isConnected {
                try? await Task.sleep(nanoseconds: 2_000_000_000)
                dataTransferred.download += Int64.random(in: 500...15000)
                dataTransferred.upload += Int64.random(in: 50...2000)
            }
        }
    }

    var formattedTime: String {
        let hours = Int(connectionTime) / 3600
        let minutes = (Int(connectionTime) % 3600) / 60
        let seconds = Int(connectionTime) % 60
        if hours > 0 {
            return String(format: "%dh %dm", hours, minutes)
        }
        return String(format: "%dm %ds", minutes, seconds)
    }

    // MARK: - Toast
    func showToast(_ message: String, type: Toast.ToastType = .info) {
        toast = Toast(message: message, type: type)
    }

    // MARK: - Deep Links
    func handleDeepLink(_ route: DeepLinkRoute) {
        switch route {
        case .connect:
            if isConnected { disconnect() } else { connect() }
        case .server(let id):
            if let server = sampleServers.first(where: { $0.id == id }) {
                selectedServer = server
                showToast("Selected \(server.city)", type: .success)
            }
        case .settings:
            showingSettings = true
        case .subscription:
            showToast("Subscription page coming soon", type: .info)
        case .promo(let code):
            showToast("Promo code applied: \(code)", type: .success)
        case .config(let url):
            showToast("Importing config from URL", type: .info)
        }
    }

    // MARK: - Error Handling
    func handleError(_ error: Error, context: String = "") {
        CrashReporter.shared.report(error, context: ["context": context])
        connectionStatus = .error(error.localizedDescription)
        showToast(error.localizedDescription, type: .error)
    }
}

enum ConnectionStatus: Equatable {
    case disconnected
    case connecting
    case connected
    case reconnecting
    case error(String)
}

struct DataUsage {
    var download: Int64 = 0
    var upload: Int64 = 0
    var formattedDownload: String {
        ByteCountFormatter.string(fromByteCount: download, countStyle: .file)
    }
    var formattedUpload: String {
        ByteCountFormatter.string(fromByteCount: upload, countStyle: .file)
    }
}

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

    var countryFlag: String {
        let base: UInt32 = 127462
        let scalars = country.uppercased().unicodeScalars
        guard scalars.count == 2 else { return "🏳️" }
        return String(String.UnicodeScalarView(
            scalars.compactMap { UnicodeScalar(base + $0.value - 65) }
        ))
    }
}

let sampleServers: [VPNServer] = [
    VPNServer(id: "auto", name: "auto.nebula", city: "Fastest Server", country: "🌍", host: "", port: 0, publicKey: "", loadPercent: 0, latencyMs: 0, isPremium: false),
    VPNServer(id: "eu-fra-01", name: "fra1.nebula", city: "Frankfurt", country: "DE", host: "fra1.nebula-vpn.io", port: 51820, publicKey: "key1", loadPercent: 42, latencyMs: 12, isPremium: false),
    VPNServer(id: "eu-ams-01", name: "ams1.nebula", city: "Amsterdam", country: "NL", host: "ams1.nebula-vpn.io", port: 51820, publicKey: "key2", loadPercent: 58, latencyMs: 18, isPremium: false),
    VPNServer(id: "gb-lon-01", name: "lon1.nebula", city: "London", country: "GB", host: "lon1.nebula-vpn.io", port: 51820, publicKey: "key6", loadPercent: 55, latencyMs: 24, isPremium: false),
    VPNServer(id: "us-nyc-01", name: "nyc1.nebula", city: "New York", country: "US", host: "nyc1.nebula-vpn.io", port: 51820, publicKey: "key3", loadPercent: 71, latencyMs: 89, isPremium: true),
    VPNServer(id: "us-mia-01", name: "mia1.nebula", city: "Miami", country: "US", host: "mia1.nebula-vpn.io", port: 51820, publicKey: "key5", loadPercent: 29, latencyMs: 130, isPremium: false),
    VPNServer(id: "us-sfo-01", name: "sfo1.nebula", city: "San Francisco", country: "US", host: "sfo1.nebula-vpn.io", port: 51820, publicKey: "key4", loadPercent: 34, latencyMs: 105, isPremium: false),
    VPNServer(id: "jp-tokyo-01", name: "tky1.nebula", city: "Tokyo", country: "JP", host: "tky1.nebula-vpn.io", port: 51820, publicKey: "key7", loadPercent: 23, latencyMs: 185, isPremium: true),
    VPNServer(id: "sg-sin-01", name: "sin1.nebula", city: "Singapore", country: "SG", host: "sin1.nebula-vpn.io", port: 51820, publicKey: "key8", loadPercent: 18, latencyMs: 210, isPremium: false),
    VPNServer(id: "au-syd-01", name: "syd1.nebula", city: "Sydney", country: "AU", host: "syd1.nebula-vpn.io", port: 51820, publicKey: "key9", loadPercent: 12, latencyMs: 280, isPremium: true),
    VPNServer(id: "br-gru-01", name: "gru1.nebula", city: "São Paulo", country: "BR", host: "gru1.nebula-vpn.io", port: 51820, publicKey: "key10", loadPercent: 9, latencyMs: 175, isPremium: false),
]
