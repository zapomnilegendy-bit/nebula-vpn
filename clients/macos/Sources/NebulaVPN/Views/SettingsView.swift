import SwiftUI
import ServiceManagement

struct SettingsView: View {
    @EnvironmentObject private var vm: MainViewModel
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        VStack(spacing: 0) {
            HStack {
                Text("Settings")
                    .font(.system(size: 15, weight: .semibold))
                Spacer()
                Button(action: { vm.saveSettings(); dismiss() }) {
                    Text("Done").font(.system(size: 13, weight: .semibold))
                }
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 14)

            Divider()

            ScrollView {
                VStack(alignment: .leading, spacing: 20) {
                    // CONNECTION
                    section("CONNECTION") {
                        Toggle("Kill Switch", isOn: $vm.settings.killSwitch)
                            .onChange(of: vm.settings.killSwitch) { _, _ in vm.saveSettings() }
                        Toggle("Adaptive Protocol", isOn: $vm.settings.adaptiveProtocol)
                        Toggle("Launch at Login", isOn: $vm.settings.launchAtLogin)
                    }

                    // PROTOCOL
                    section("PROTOCOL") {
                        Picker("Default Protocol", selection: $vm.settings.selectedProtocol) {
                            ForEach(TransportProtocol.allCases) { p in
                                Text(p.displayName).tag(p)
                            }
                        }
                        .pickerStyle(.menu)
                    }

                    // DNS
                    section("DNS") {
                        Toggle("DNS over HTTPS", isOn: $vm.settings.dnsOverHTTPS)
                        Toggle("DNS over TLS", isOn: $vm.settings.dnsOverTLS)
                    }

                    // APPEARANCE
                    section("APPEARANCE") {
                        Picker("Theme", selection: $vm.settings.appearance) {
                            Text("System").tag(Appearance.system)
                            Text("Light").tag(Appearance.light)
                            Text("Dark").tag(Appearance.dark)
                        }
                        .pickerStyle(.segmented)
                    }

                    // ACCOUNT
                    section("ACCOUNT") {
                        HStack {
                            Text("Plan")
                            Spacer()
                            Text("Pro")
                                .foregroundColor(.secondary)
                        }
                        HStack {
                            Text("Devices")
                            Spacer()
                            Text("2 of 5")
                                .foregroundColor(.secondary)
                        }
                    }

                    // ABOUT
                    section("ABOUT") {
                        HStack {
                            Text("Version")
                            Spacer()
                            Text("1.0.0 (1)")
                                .foregroundColor(.secondary)
                        }
                        Button("Check for Updates") {
                            AutoUpdater().checkForUpdates()
                        }
                        .buttonStyle(.plain)
                        .foregroundColor(.accentColor)
                        .font(.system(size: 13))
                    }
                }
                .padding(20)
            }
        }
        .frame(width: 360, height: 480)
        .background(VisualEffectView(material: .sidebar, blendingMode: .behindWindow))
    }

    private func section(_ title: String, @ViewBuilder content: () -> some View) -> some View {
        VStack(alignment: .leading, spacing: 8) {
            Text(title)
                .font(.system(size: 10, weight: .semibold))
                .foregroundColor(.secondary)
                .kerning(0.8)
            VStack(spacing: 0) {
                content()
                    .font(.system(size: 13))
            }
            .padding(12)
            .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 10))
        }
    }
}
