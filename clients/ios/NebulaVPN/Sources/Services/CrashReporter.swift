import UIKit
import OSLog

final class CrashReporter {
    private let apiURL = URL(string: "https://api.nebula-vpn.io/v1/crashes")!
    private let log = Logger(subsystem: "com.nebula.vpn.ios", category: "CrashReporter")
    private let deviceID: String

    static let shared = CrashReporter()

    private init() {
        deviceID = "ios-" + (UIDevice.current.identifierForVendor?.uuidString.prefix(8) ?? UUID().uuidString.prefix(8)).lowercased()
    }

    func start() {
        setupExceptionHandler()
        setupSignalHandlers()
        log.info("CrashReporter initialized")
    }

    private func setupExceptionHandler() {
        NSSetUncaughtExceptionHandler { exception in
            let info: [String: Any] = [
                "name": exception.name.rawValue,
                "reason": exception.reason ?? "",
                "callStack": exception.callStackSymbols,
                "platform": "ios",
                "version": Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as? String ?? "1.0.0",
                "device": UIDevice.current.model,
                "system": UIDevice.current.systemVersion,
            ]
            Self.sendCrash(info)
        }
    }

    private func setupSignalHandlers() {
        let signals = [SIGABRT, SIGSEGV, SIGBUS, SIGILL, SIGFPE, SIGPIPE]
        signals.forEach { sig in
            signal(sig) { _ in
                var info: [String: Any] = ["signal": Int(sig), "platform": "ios"]
                Self.sendCrash(info)
                exit(sig)
            }
        }
    }

    func report(_ error: Error, context: [String: Any] = [:]) {
        var payload: [String: Any] = [
            "error": error.localizedDescription,
            "platform": "ios",
            "version": Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as? String ?? "1.0.0",
            "device_id": deviceID,
        ]
        payload.merge(context) { $1 }
        sendToAPI(payload)
    }

    func reportMessage(_ msg: String, level: String = "error", context: [String: Any] = [:]) {
        var payload: [String: Any] = ["message": msg, "level": level, "platform": "ios", "device_id": deviceID]
        payload.merge(context) { $1 }
        sendToAPI(payload)
    }

    private static func sendCrash(_ info: [String: Any]) {
        guard let data = try? JSONSerialization.data(withJSONObject: info),
              let url = URL(string: "https://api.nebula-vpn.io/v1/crashes") else { return }
        var req = URLRequest(url: url)
        req.httpMethod = "POST"
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = data
        _ = try? URLSession.shared.synchronousUpload(req)
    }

    private func sendToAPI(_ payload: [String: Any]) {
        Task {
            guard let data = try? JSONSerialization.data(withJSONObject: payload) else { return }
            var req = URLRequest(url: apiURL)
            req.httpMethod = "POST"
            req.setValue("application/json", forHTTPHeaderField: "Content-Type")
            req.httpBody = data
            _ = try? await URLSession.shared.data(for: req)
        }
    }
}

private extension URLSession {
    func synchronousUpload(_ request: URLRequest) throws -> (Data, URLResponse) {
        let semaphore = DispatchSemaphore(value: 0)
        var result: (Data, URLResponse)!
        let task = self.dataTask(with: request) { data, response, error in
            result = (data ?? Data(), response ?? URLResponse())
            semaphore.signal()
        }
        task.resume()
        semaphore.wait()
        return result
    }
}
