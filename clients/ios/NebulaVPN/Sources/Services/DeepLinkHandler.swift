import UIKit
import OSLog

enum DeepLinkRoute: Equatable {
    case connect
    case server(id: String)
    case settings
    case subscription
    case promo(code: String)
    case config(url: URL)
}

final class DeepLinkHandler {
    private let log = Logger(subsystem: "com.nebula.vpn.ios", category: "DeepLink")

    static let shared = DeepLinkHandler()

    private init() {}

    func handle(_ url: URL) -> DeepLinkRoute? {
        log.info("Handling URL: \(url.absoluteString)")

        guard let components = URLComponents(url: url, resolvingAgainstBaseURL: false) else {
            log.error("Failed to parse URL components")
            return nil
        }

        switch components.scheme?.lowercased() {
        case "nebula":
            return handleNebulaScheme(components)
        case "https", "http":
            return handleUniversalLink(components)
        default:
            log.warning("Unknown scheme: \(components.scheme ?? "nil")")
            return nil
        }
    }

    private func handleNebulaScheme(_ components: URLComponents) -> DeepLinkRoute? {
        let host = components.host?.lowercased()

        switch host {
        case "connect":
            return .connect
        case "server":
            if let id = components.queryItems?.first(where: { $0.name == "id" })?.value {
                return .server(id: id)
            }
            return nil
        case "settings":
            return .settings
        case "subscription":
            return .subscription
        case "promo":
            if let code = components.queryItems?.first(where: { $0.name == "code" })?.value {
                return .promo(code: code)
            }
            return nil
        case "import":
            if let urlStr = components.queryItems?.first(where: { $0.name == "url" })?.value,
               let url = URL(string: urlStr) {
                return .config(url: url)
            }
            return nil
        default:
            log.warning("Unknown host: \(host ?? "nil")")
            return nil
        }
    }

    private func handleUniversalLink(_ components: URLComponents) -> DeepLinkRoute? {
        let path = components.path.lowercased()

        switch path {
        case "/connect":
            return .connect
        case let p where p.hasPrefix("/server/"):
            let id = String(p.dropFirst("/server/".count))
            return .server(id: id)
        case "/settings":
            return .settings
        case "/subscription":
            return .subscription
        case let p where p.hasPrefix("/promo/"):
            let code = String(p.dropFirst("/promo/".count))
            return .promo(code: code)
        default:
            return nil
        }
    }
}

// MARK: - Scene Delegate Integration
extension NebulaVPNApp {
    func application(_ app: UIApplication, open url: URL, options: [UIApplication.OpenURLOptionsKey: Any] = [:]) -> Bool {
        return DeepLinkHandler.shared.handle(url) != nil
    }

    func application(_ application: UIApplication, continue userActivity: NSUserActivity, restorationHandler: @escaping ([UIUserActivityRestoring]) -> Void) -> Bool {
        guard userActivity.activityType == NSUserActivityTypeBrowsingWeb,
              let url = userActivity.webpageURL else { return false }
        return DeepLinkHandler.shared.handle(url) != nil
    }
}
