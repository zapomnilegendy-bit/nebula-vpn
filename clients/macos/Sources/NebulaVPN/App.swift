import SwiftUI

@main
struct NebulaVPNApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    @StateObject private var vm = MainViewModel()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(vm)
                .frame(minWidth: 380, maxWidth: 420, minHeight: 520, maxHeight: 600)
                .onAppear {
                    NSWindow.allowsAutomaticWindowTabbing = false
                }
        }
        .windowStyle(.hiddenTitleBar)
        .windowResizability(.contentSize)
        .commands {
            CommandGroup(replacing: .appInfo) {
                Button("About NebulaVPN") { vm.showingAbout = true }
            }
        }
    }
}

class AppDelegate: NSObject, NSApplicationDelegate {
    private var crashReporter: CrashReporter?
    private var updater: AutoUpdater?

    func applicationDidFinishLaunching(_ notification: Notification) {
        crashReporter = CrashReporter()
        updater = AutoUpdater()
        updater?.checkForUpdates()

        DispatchQueue.main.asyncAfter(deadline: .now() + 30) {
            self.updater?.checkForUpdatesInBackground()
        }
    }
}
