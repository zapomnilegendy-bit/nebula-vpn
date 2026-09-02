import Foundation
import Security

final class KeychainManager {
    private let service = "com.nebula.vpn"
    private let settingsKey = "app_settings"
    private let onboardingKey = "onboarding_complete"

    // MARK: - Settings
    func saveSettings(_ settings: AppSettings) {
        guard let data = try? JSONEncoder().encode(settings) else { return }
        save(data: data, key: settingsKey)
    }

    func loadSettings() -> AppSettings? {
        guard let data = load(key: settingsKey) else { return nil }
        return try? JSONDecoder().decode(AppSettings.self, from: data)
    }

    // MARK: - Onboarding
    func markOnboardingSeen() {
        save(data: Data([1]), key: onboardingKey)
    }

    func hasSeenOnboarding() -> Bool {
        load(key: onboardingKey) != nil
    }

    // MARK: - Secrets
    func saveAuthToken(_ token: String) {
        save(data: Data(token.utf8), key: "auth_token")
    }

    func loadAuthToken() -> String? {
        guard let data = load(key: "auth_token") else { return nil }
        return String(data: data, encoding: .utf8)
    }

    func savePrivateKey(_ key: String) {
        save(data: Data(key.utf8), key: "wg_private_key")
    }

    func loadPrivateKey() -> String? {
        guard let data = load(key: "wg_private_key") else { return nil }
        return String(data: data, encoding: .utf8)
    }

    // MARK: - Keychain Operations
    private func save(data: Data, key: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecValueData as String: data,
            kSecAttrAccessible as String: kSecAttrAccessibleWhenUnlockedThisDeviceOnly,
        ]

        SecItemDelete(query as CFDictionary)
        SecItemAdd(query as CFDictionary, nil)
    }

    private func load(key: String) -> Data? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne,
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        guard status == errSecSuccess else { return nil }
        return result as? Data
    }

    func clearAll() {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
        ]
        SecItemDelete(query as CFDictionary)
    }
}
