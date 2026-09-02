import SwiftUI

struct ServerListView: View {
    @EnvironmentObject private var appState: AppState
    @Environment(\.dismiss) private var dismiss
    @State private var searchText = ""

    private let servers = sampleServers

    var filteredServers: [VPNServer] {
        if searchText.isEmpty { return servers }
        return servers.filter {
            $0.city.localizedCaseInsensitiveContains(searchText) ||
            $0.country.localizedCaseInsensitiveContains(searchText) ||
            $0.name.localizedCaseInsensitiveContains(searchText)
        }
    }

    var body: some View {
        NavigationStack {
            List {
                // Fastest server
                Section {
                    Button(action: {
                        selectFastest()
                    }) {
                        HStack(spacing: 14) {
                            ZStack {
                                RoundedRectangle(cornerRadius: 10)
                                    .fill(Color.accentColor.opacity(0.1))
                                    .frame(width: 40, height: 40)
                                Image(systemName: "bolt.fill")
                                    .font(.system(size: 16))
                                    .foregroundColor(.accentColor)
                            }
                            VStack(alignment: .leading, spacing: 2) {
                                Text("Fastest Server")
                                    .font(.system(size: 15, weight: .semibold))
                                Text("Automatically selected")
                                    .font(.system(size: 12))
                                    .foregroundColor(.secondary)
                            }
                            .foregroundColor(.primary)
                            Spacer()
                            if appState.selectedServer?.id == "auto" {
                                Image(systemName: "checkmark")
                                    .font(.system(size: 14, weight: .semibold))
                                    .foregroundColor(.accentColor)
                            }
                        }
                    }
                }

                // All locations
                Section("All Locations") {
                    ForEach(filteredServers.filter { $0.id != "auto" }) { server in
                        Button(action: { select(server) }) {
                            ServerRow(server: server)
                                .foregroundColor(.primary)
                        }
                    }
                }
            }
            .searchable(text: $searchText, prompt: "Search locations")
            .navigationTitle("Servers")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Done") { dismiss() }
                        .font(.system(size: 15, weight: .semibold))
                }
            }
        }
    }

    private func select(_ server: VPNServer) {
        appState.selectedServer = server
        dismiss()
    }

    private func selectFastest() {
        let fastest = servers
            .filter { $0.id != "auto" && $0.latencyMs > 0 }
            .min(by: { $0.latencyMs < $1.latencyMs })
        appState.selectedServer = fastest
        dismiss()
    }
}

struct ServerRow: View {
    let server: VPNServer
    @EnvironmentObject private var appState: AppState

    var body: some View {
        HStack(spacing: 14) {
            Text(server.countryFlag)
                .font(.title3)
            VStack(alignment: .leading, spacing: 2) {
                Text(server.city)
                    .font(.system(size: 15, weight: .semibold))
                HStack(spacing: 4) {
                    Text(server.name)
                        .font(.system(size: 12, design: .monospaced))
                        .foregroundColor(.secondary)
                    if server.isPremium {
                        Text("PREMIUM")
                            .font(.system(size: 8, weight: .bold))
                            .foregroundColor(.white)
                            .padding(.horizontal, 4)
                            .padding(.vertical, 1)
                            .background(Color.orange, in: RoundedRectangle(cornerRadius: 3))
                    }
                }
            }
            Spacer()
            VStack(alignment: .trailing, spacing: 4) {
                HStack(spacing: 4) {
                    Circle()
                        .fill(server.loadPercent < 60 ? Color.green : server.loadPercent < 85 ? Color.orange : Color.red)
                        .frame(width: 6, height: 6)
                    Text("\(server.latencyMs)ms")
                        .font(.system(size: 13, weight: .medium, design: .monospaced))
                        .foregroundColor(.secondary)
                }
            }
        }
        .padding(.vertical, 4)
    }
}


