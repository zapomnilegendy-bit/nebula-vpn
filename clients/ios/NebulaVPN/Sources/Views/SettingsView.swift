import SwiftUI

struct SettingsView: View {
    @EnvironmentObject private var appState: AppState
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        NavigationStack {
            Form {
                Section("Connection") {
                    Toggle("Kill Switch", isOn: .constant(true))
                    Toggle("Split Tunneling", isOn: .constant(false))
                    Toggle("Auto-protocol", isOn: .constant(true))
                }

                Section("DNS") {
                    NavigationLink("DNS Configuration") {
                        List {
                            Toggle("DNS over HTTPS", isOn: .constant(true))
                            Toggle("DNS over TLS", isOn: .constant(true))
                            Toggle("Prevent DNS leaks", isOn: .constant(true))
                        }
                        .navigationTitle("DNS")
                    }
                }

                Section("Account") {
                    HStack {
                        Text("Plan")
                        Spacer()
                        Text("Pro")
                            .foregroundColor(.secondary)
                    }
                    HStack {
                        Text("Data Used")
                        Spacer()
                        Text("2.4 GB")
                            .foregroundColor(.secondary)
                    }
                }

                Section("About") {
                    HStack {
                        Text("Version")
                        Spacer()
                        Text("1.0.0")
                            .foregroundColor(.secondary)
                    }
                }
            }
            .navigationTitle("Settings")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .confirmationAction) {
                    Button("Done") { dismiss() }
                        .font(.system(size: 15, weight: .semibold))
                }
            }
        }
    }
}
