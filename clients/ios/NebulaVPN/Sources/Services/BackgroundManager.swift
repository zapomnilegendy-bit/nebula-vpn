import Foundation
import UIKit
import BackgroundTasks

/// Manages background tasks with battery-efficient scheduling
/// Uses BGTaskScheduler for periodic health checks and reconnections
@MainActor
final class BackgroundManager {
    static let shared = BackgroundManager()

    private let healthCheckTaskId = "com.nebula.vpn.healthcheck"
    private let reconnectTaskId = "com.nebula.vpn.reconnect"
    private let minimumInterval: TimeInterval = 4 * 60 * 60 // 4 hours minimum

    private var isRegistered = false

    private init() {
        registerTasks()
    }

    // MARK: - Registration

    func registerTasks() {
        guard !isRegistered else { return }
        BGTaskScheduler.shared.register(forTaskWithIdentifier: healthCheckTaskId, using: nil) { [weak self] task in
            self?.handleHealthCheck(task as! BGProcessingTask)
        }
        BGTaskScheduler.shared.register(forTaskWithIdentifier: reconnectTaskId, using: nil) { [weak self] task in
            self?.handleReconnect(task as! BGAppRefreshTask)
        }
        isRegistered = true
    }

    func scheduleTasks() {
        scheduleHealthCheck()
        scheduleReconnect()
    }

    // MARK: - Health Check

    private func scheduleHealthCheck() {
        let request = BGProcessingTaskRequest(identifier: healthCheckTaskId)
        request.requiresNetworkConnectivity = true
        request.requiresExternalPower = false
        request.earliestBeginDate = Date(timeIntervalSinceNow: minimumInterval)

        do {
            try BGTaskScheduler.shared.submit(request)
        } catch {
            os_log(.error, "Failed to schedule health check: %@", error.localizedDescription)
        }
    }

    private func handleHealthCheck(_ task: BGProcessingTask) {
        task.expirationHandler = { task.setTaskCompleted(success: false) }

        Task {
            do {
                try await performHealthCheck()
                task.setTaskCompleted(success: true)
            } catch {
                task.setTaskCompleted(success: false)
            }
            scheduleHealthCheck()
        }
    }

    private func performHealthCheck() async throws {
        let api = VpnAPI()
        let health = try await api.checkHealth()

        if !health.isHealthy {
            // Trigger reconnect if needed
            NotificationCenter.default.post(
                name: .nebulaNeedsReconnect,
                object: nil
            )
        }
    }

    // MARK: - Reconnect

    private func scheduleReconnect() {
        let request = BGAppRefreshTaskRequest(identifier: reconnectTaskId)
        request.earliestBeginDate = Date(timeIntervalSinceNow: 3600) // 1 hour

        do {
            try BGTaskScheduler.shared.submit(request)
        } catch {
            os_log(.error, "Failed to schedule reconnect: %@", error.localizedDescription)
        }
    }

    private func handleReconnect(_ task: BGAppRefreshTask) {
        task.expirationHandler = { task.setTaskCompleted(success: false) }

        Task {
            do {
                try await attemptReconnect()
                task.setTaskCompleted(success: true)
            } catch {
                task.setTaskCompleted(success: false)
            }
            scheduleReconnect()
        }
    }

    private func attemptReconnect() async throws {
        guard let session = try? await VpnManager.shared.getStatus(),
              session.status == .disconnected else { return }

        // Smart reconnect with best server
        try await VpnManager.shared.reconnect()
    }

    // MARK: - Push Registration

    func registerForPushNotifications() async {
        do {
            let granted = try await UNUserNotificationCenter.current().requestAuthorization(
                options: [.alert, .sound, .badge]
            )
            guard granted else { return }

            let settings = await UNUserNotificationCenter.current().notificationSettings()
            guard settings.authorizationStatus == .authorized else { return }

            #if !targetEnvironment(simulator)
            let token = await UIApplication.shared.registerForRemoteNotifications()
            // Register token with backend
            try await VpnAPI.shared.registerPushToken(token)
            #endif
        } catch {
            os_log(.error, "Push registration failed: %@", error.localizedDescription)
        }
    }
}

// MARK: - Notifications

extension Notification.Name {
    static let nebulaNeedsReconnect = Notification.Name("com.nebula.vpn.needsReconnect")
    static let nebulaBatteryLow = Notification.Name("com.nebula.vpn.batteryLow")
}
