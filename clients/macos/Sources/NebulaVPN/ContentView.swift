import SwiftUI

struct ContentView: View {
    @EnvironmentObject private var vm: MainViewModel

    var body: some View {
        ZStack {
            VisualEffectView(material: .sidebar, blendingMode: .behindWindow)
                .ignoresSafeArea()

            if vm.isOnboarding {
                OnboardingView()
            } else if !vm.isAuthenticated {
                LockScreenView()
            } else {
                MainView()
            }
        }
        .sheet(isPresented: $vm.showingSettings) { SettingsView() }
        .sheet(isPresented: $vm.showingAbout) { AboutView() }
    }
}

// MARK: - Lock Screen
struct LockScreenView: View {
    @EnvironmentObject private var vm: MainViewModel

    var body: some View {
        VStack(spacing: 20) {
            Image(systemName: "lock.shield.fill")
                .font(.system(size: 48))
                .foregroundColor(.secondary)
            Text("NebulaVPN")
                .font(.system(size: 24, weight: .bold))
            Text("Authenticate to access your secure connection")
                .font(.system(size: 13))
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
            Button("Unlock") {
                vm.authenticateUser()
            }
            .buttonStyle(.borderedProminent)
        }
        .padding(40)
    }
}

// MARK: - Main View
struct MainView: View {
    @EnvironmentObject private var vm: MainViewModel
    @State private var pulsing = false

    var body: some View {
        VStack(spacing: 0) {
            // Top bar
            HStack {
                trafficBar
                Spacer()
                Button(action: { vm.showingSettings = true }) {
                    Image(systemName: "gearshape")
                        .font(.system(size: 13, weight: .medium))
                        .foregroundColor(.secondary)
                        .frame(width: 28, height: 28)
                        .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 7))
                }
                .buttonStyle(.plain)
            }
            .padding(.horizontal, 20)
            .padding(.top, 12)

            Spacer()

            // Connection ring
            ConnectionRingView(
                status: vm.connectionStatus,
                progress: vm.isConnected ? 0.82 : 0.18
            )
            .frame(width: 200, height: 200)
            .padding(.bottom, 28)

            // Status
            VStack(spacing: 4) {
                Text(vm.connectionStatus.label)
                    .font(.system(size: 16, weight: .semibold))
                if vm.isConnected {
                    Text(vm.formattedDuration)
                        .font(.system(size: 13, design: .monospaced))
                        .foregroundColor(.secondary)
                }
            }

            // Data usage
            if vm.isConnected {
                HStack(spacing: 24) {
                    dataBadge(icon: "arrow.down", value: vm.dataTransferred.formattedDownload)
                    dataBadge(icon: "arrow.up", value: vm.dataTransferred.formattedUpload)
                }
                .padding(.top, 16)
                .transition(.scale.combined(with: .opacity))
            }

            Spacer()

            // Server selector
            Button(action: { vm.showingServerPicker = true }) {
                HStack(spacing: 12) {
                    if let server = vm.selectedServer {
                        Text(server.countryFlag)
                            .font(.title2)
                        VStack(alignment: .leading, spacing: 1) {
                            Text(server.city)
                                .font(.system(size: 14, weight: .semibold))
                                .foregroundColor(.primary)
                            Text(server.detailString)
                                .font(.system(size: 11))
                                .foregroundColor(.secondary)
                        }
                    } else {
                        Image(systemName: "magnifyingglass")
                            .foregroundColor(.secondary)
                        Text("Select a server")
                            .font(.system(size: 14, weight: .medium))
                            .foregroundColor(.secondary)
                    }
                    Spacer()
                    Image(systemName: "chevron.forward")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 12)
                .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 12))
                .overlay(RoundedRectangle(cornerRadius: 12).stroke(Color.primary.opacity(0.05), lineWidth: 1))
            }
            .buttonStyle(.plain)
            .padding(.horizontal, 20)
            .padding(.bottom, 12)

            // Connect button
            Button(action: { vm.isConnected ? vm.disconnect() : vm.connect() }) {
                Text(vm.isConnected ? "Disconnect" : "Connect")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundColor(vm.isConnected ? .red : .white)
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 12)
                    .background(
                        vm.isConnected
                        ? Color.red.opacity(0.1)
                        : Color.primary,
                        in: RoundedRectangle(cornerRadius: 12)
                    )
            }
            .buttonStyle(.plain)
            .padding(.horizontal, 20)
            .padding(.bottom, 20)
        }
    }

    private var trafficBar: some View {
        HStack(spacing: 0) {
            Image(systemName: "arrow.down.circle.fill")
                .font(.caption)
                .foregroundColor(.secondary)
            Text(vm.dataTransferred.formattedDownload)
                .font(.system(size: 11, weight: .medium, design: .monospaced))
                .foregroundColor(.secondary)
                .frame(width: 60, alignment: .leading)
            Spacer()
            Image(systemName: "arrow.up.circle.fill")
                .font(.caption)
                .foregroundColor(.secondary)
            Text(vm.dataTransferred.formattedUpload)
                .font(.system(size: 11, weight: .medium, design: .monospaced))
                .foregroundColor(.secondary)
                .frame(width: 60, alignment: .leading)
        }
        .padding(.vertical, 6)
        .padding(.horizontal, 12)
        .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 8))
        .opacity(vm.isConnected ? 1 : 0)
        .animation(.easeInOut(duration: 0.3), value: vm.isConnected)
    }

    private func dataBadge(icon: String, value: String) -> some View {
        HStack(spacing: 5) {
            Image(systemName: icon)
                .font(.system(size: 10, weight: .semibold))
                .foregroundColor(.secondary)
            Text(value)
                .font(.system(size: 12, weight: .medium, design: .monospaced))
                .foregroundColor(.secondary)
        }
        .padding(.horizontal, 12)
        .padding(.vertical, 6)
        .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 8))
    }
}

struct VisualEffectView: NSViewRepresentable {
    let material: NSVisualEffectView.Material
    let blendingMode: NSVisualEffectView.BlendingMode

    func makeNSView(context: Context) -> NSVisualEffectView {
        let v = NSVisualEffectView()
        v.material = material
        v.blendingMode = blendingMode
        v.state = .active
        return v
    }
    func updateNSView(_ v: NSVisualEffectView, context: Context) {
        v.material = material; v.blendingMode = blendingMode
    }
}
