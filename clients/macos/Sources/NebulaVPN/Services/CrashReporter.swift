import Foundation

final class CrashReporter {
    private let apiURL = URL(string: "https://api.nebula-vpn.io/v1/crashes")!
    private let userID: String

    init() {
        userID = "macos-" + UUID().uuidString.prefix(8).lowercased()
        setupCrashHandling()
    }

    private func setupCrashHandling() {
        NSSetUncaughtExceptionHandler { exception in
            let info: [String: Any] = [
                "name": exception.name.rawValue,
                "reason": exception.reason ?? "",
                "callStack": exception.callStackSymbols,
                "platform": "macos",
                "version": "1.0.0",
            ]
            if let data = try? JSONSerialization.data(withJSONObject: info) {
                var request = URLRequest(url: URL(string: "https://api.nebula-vpn.io/v1/crashes")!)
                request.httpMethod = "POST"
                request.setValue("application/json", forHTTPHeaderField: "Content-Type")
                request.httpBody = data
                URLSession.shared.synchronousDataTask(with: request)
            }
        }

        signal(SIGABRT, crashSignalHandler)
        signal(SIGSEGV, crashSignalHandler)
        signal(SIGBUS, crashSignalHandler)
        signal(SIGILL, crashSignalHandler)
    }

    func reportError(_ error: Error, context: [String: Any] = [:]) {
        var payload: [String: Any] = [
            "error": error.localizedDescription,
            "platform": "macos",
            "version": "1.0.0",
            "user_id": userID,
        ]
        payload.merge(context) { $1 }

        Task {
            guard let data = try? JSONSerialization.data(withJSONObject: payload) else { return }
            var request = URLRequest(url: apiURL)
            request.httpMethod = "POST"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = data
            _ = try? await URLSession.shared.data(for: request)
        }
    }
}

private let crashReportURL = URL(string: "https://api.nebula-vpn.io/v1/crashes")!

private func crashSignalHandler(_ signal: Int32) {
    let info: [String: Any] = [
        "signal": Int(signal),
        "platform": "macos",
        "version": "1.0.0",
    ]
    if let data = try? JSONSerialization.data(withJSONObject: info) {
        var request = URLRequest(url: crashReportURL)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = data
        URLSession.shared.synchronousDataTask(with: request)
    }
    exit(signal)
}

extension URLSession {
    func synchronousDataTask(with request: URLRequest) {
        let semaphore = DispatchSemaphore(value: 0)
        let task = self.dataTask(with: request) { _, _, _ in semaphore.signal() }
        task.resume()
        _ = semaphore.wait(timeout: .now() + 5)
    }
}
