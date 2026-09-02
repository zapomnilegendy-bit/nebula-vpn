import Foundation
import UIKit

/// Battery-optimized background task manager
/// Uses deferred processing, adaptive polling, and app state observation
final class BatteryManager {
    static let shared = BatteryManager()

    private var pollTimer: Timer?
    private var isAppActive = true
    private var idleCount = 0
    private var batchBuffer: [Data] = []
    private let batchLimit = 64 * 1024 // 64KB batch flush
    private let processQueue = DispatchQueue(
        label: "com.nebula.vpn.battery",
        qos: .utility,
        autoreleaseFrequency: .workItem
    )

    // Adaptive poll intervals (ms)
    private let activeInterval: TimeInterval = 0.05  // 50ms when active
    private let idleInterval: TimeInterval = 0.2     // 200ms when idle
    private let backgroundInterval: TimeInterval = 2.0 // 2s in background

    private init() {
        setupObservers()
    }

    // MARK: - Setup

    private func setupObservers() {
        NotificationCenter.default.addObserver(
            self,
            selector: #selector(appDidBecomeActive),
            name: UIApplication.didBecomeActiveNotification,
            object: nil
        )
        NotificationCenter.default.addObserver(
            self,
            selector: #selector(appDidEnterBackground),
            name: UIApplication.didEnterBackgroundNotification,
            object: nil
        )
        NotificationCenter.default.addObserver(
            self,
            selector: #selector(appWillResignActive),
            name: UIApplication.willResignActiveNotification,
            object: nil
        )
    }

    // MARK: - Lifecycle

    func start() {
        DispatchQueue.main.async { [weak self] in
            self?.startPolling(interval: self?.activeInterval ?? 0.05)
        }
    }

    func stop() {
        DispatchQueue.main.async { [weak self] in
            self?.pollTimer?.invalidate()
            self?.pollTimer = nil
        }
        flushBatch()
    }

    // MARK: - Deferred Processing

    func enqueue(_ data: Data) {
        processQueue.async { [weak self] in
            guard let self = self else { return }
            self.batchBuffer.append(data)

            if self.batchBuffer.reduce(0, { $0 + $1.count }) >= self.batchLimit {
                self.flushBatch()
            }
        }
    }

    func enqueueDelayed(_ block: @escaping () -> Void) {
        processQueue.asyncAfter(deadline: .now() + 0.05) {
            block()
        }
    }

    private func flushBatch() {
        guard !batchBuffer.isEmpty else { return }
        let batch = batchBuffer
        batchBuffer.removeAll(keepingCapacity: true)

        // Process batch (e.g., send to packet tunnel)
        processBatch(batch)
    }

    private func processBatch(_ batch: [Data]) {
        // Delegate to tunnel provider
        TunnelMessageBus.shared.send(batch)
    }

    // MARK: - Adaptive Poll Interval

    private func startPolling(interval: TimeInterval) {
        pollTimer?.invalidate()
        pollTimer = Timer.scheduledTimer(
            timeInterval: interval,
            target: self,
            selector: #selector(pollTick),
            userInfo: nil,
            repeats: true
        )
    }

    @objc private func pollTick() {
        guard isAppActive else {
            // Defer non-critical work
            idleCount += 1
            if idleCount > 10 {
                adjustPollInterval(isIdle: true)
            }
            return
        }
        idleCount = 0
        adjustPollInterval(isIdle: false)
    }

    private func adjustPollInterval(isIdle: Bool) {
        let newInterval: TimeInterval
        if !isAppActive {
            newInterval = backgroundInterval
        } else if isIdle {
            newInterval = idleInterval
        } else {
            newInterval = activeInterval
        }

        DispatchQueue.main.async { [weak self] in
            guard let self = self,
                  let timer = self.pollTimer,
                  abs(timer.timeInterval - newInterval) > 0.01 else { return }
            self.startPolling(interval: newInterval)
        }
    }

    // MARK: - App State

    @objc private func appDidBecomeActive() {
        isAppActive = true
        idleCount = 0
        startPolling(interval: activeInterval)
    }

    @objc private func appWillResignActive() {
        isAppActive = false
    }

    @objc private func appDidEnterBackground() {
        isAppActive = false
        startPolling(interval: backgroundInterval)
        flushBatch()
    }

    deinit {
        NotificationCenter.default.removeObserver(self)
        stop()
    }
}

// Simple message bus for tunnel communication
final class TunnelMessageBus {
    static let shared = TunnelMessageBus()

    private let queue = DispatchQueue(label: "com.nebula.vpn.tunnel-bus", qos: .utility)

    private var handlers: [(Data) -> Void] = []

    private init() {}

    func subscribe(_ handler: @escaping (Data) -> Void) {
        queue.async { [weak self] in
            self?.handlers.append(handler)
        }
    }

    func send(_ data: Data) {
        queue.async { [weak self] in
            guard let self = self else { return }
            for handler in self.handlers {
                handler(data)
            }
        }
    }

    func send(_ batch: [Data]) {
        queue.async { [weak self] in
            guard let self = self else { return }
            for data in batch {
                for handler in self.handlers {
                    handler(data)
                }
            }
        }
    }
}
