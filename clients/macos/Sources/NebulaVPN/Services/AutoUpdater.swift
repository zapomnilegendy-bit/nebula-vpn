import Foundation
import AppKit

final class AutoUpdater {
    private let updateURL = URL(string: "https://api.nebula-vpn.io/v1/updates/macos")!
    private let currentVersion = "1.0.0"
    private var isChecking = false

    func checkForUpdates() {
        isChecking = true
        Task {
            defer { isChecking = false }
            guard let update = try? await fetchUpdateInfo() else { return }

            if update.version.compare(currentVersion, options: .numeric) == .orderedDescending {
                await showUpdateAlert(update)
            }
        }
    }

    func checkForUpdatesInBackground() {
        Task {
            guard !isChecking else { return }
            guard let update = try? await fetchUpdateInfo() else { return }

            if update.version.compare(currentVersion, options: .numeric) == .orderedDescending {
                await showUpdateAlert(update)
            }
        }
    }

    private func fetchUpdateInfo() async throws -> UpdateInfo {
        var request = URLRequest(url: updateURL, cachePolicy: .reloadIgnoringCacheData)
        request.timeoutInterval = 10

        let (data, _) = try await URLSession.shared.data(for: request)
        return try JSONDecoder().decode(UpdateInfo.self, from: data)
    }

    @MainActor
    private func showUpdateAlert(_ update: UpdateInfo) {
        let alert = NSAlert()
        alert.messageText = "Update Available"
        alert.informativeText = "NebulaVPN \(update.version) is available. You have \(currentVersion)."
        alert.alertStyle = .informational
        alert.addButton(withTitle: "Download")
        alert.addButton(withTitle: "Later")

        if alert.runModal() == .alertFirstButtonReturn, let url = URL(string: update.downloadURL) {
            NSWorkspace.shared.open(url)
        }
    }
}

struct UpdateInfo: Codable {
    let version: String
    let downloadURL: String
    let releaseNotes: String
    let minimumOSVersion: String
}
