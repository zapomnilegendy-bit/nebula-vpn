import SwiftUI

struct ServerPickerView: View {
    @EnvironmentObject private var vm: MainViewModel
    @Environment(\.dismiss) private var dismiss
    @State private var searchText = ""

    private var filtered: [VPNServer] {
        if searchText.isEmpty { return vm.servers }
        return vm.servers.filter {
            $0.city.localizedCaseInsensitiveContains(searchText) ||
            $0.country.localizedCaseInsensitiveContains(searchText) ||
            $0.name.localizedCaseInsensitiveContains(searchText)
        }
    }

    var body: some View {
        VStack(spacing: 0) {
            // Header
            HStack {
                Text("Servers")
                    .font(.system(size: 15, weight: .semibold))
                Spacer()
                HStack(spacing: 4) {
                    Button(action: { vm.showingQRImporter = true }) {
                        Image(systemName: "qrcode.viewfinder")
                            .font(.system(size: 13))
                    }
                    Button(action: { vm.selectFastest() }) {
                        Image(systemName: "bolt")
                            .font(.system(size: 13))
                    }
                    Button(action: { dismiss() }) {
                        Image(systemName: "xmark.circle.fill")
                            .font(.system(size: 18))
                            .foregroundColor(.secondary)
                    }
                }
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 14)

            Divider()

            // Search
            HStack {
                Image(systemName: "magnifyingglass")
                    .font(.system(size: 12))
                    .foregroundColor(.secondary)
                TextField("Search locations…", text: $searchText)
                    .textFieldStyle(.plain)
                    .font(.system(size: 13))
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 8)
            .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 8))
            .padding(.horizontal, 16)
            .padding(.vertical, 10)

            // List
            ScrollView {
                VStack(spacing: 2) {
                    SectionHeader("RECOMMENDED")
                    serverRow(VPNServer(
                        id: "auto", name: "auto.nebula", city: "Fastest Server",
                        country: "🌍", host: "", port: 0, publicKey: "",
                        loadPercent: 0, latencyMs: 0, isPremium: false,
                        supportedProtocols: [.wireGuard]
                    ))

                    SectionHeader("ALL LOCATIONS")
                    ForEach(filtered) { server in
                        serverRow(server)
                    }
                }
                .padding(.horizontal, 12)
                .padding(.bottom, 12)
            }
        }
        .frame(width: 340, height: 460)
        .background(VisualEffectView(material: .sidebar, blendingMode: .behindWindow))
    }

    private func SectionHeader(_ title: String) -> some View {
        HStack {
            Text(title)
                .font(.system(size: 10, weight: .semibold))
                .foregroundColor(.secondary)
                .kerning(0.8)
            Spacer()
        }
        .padding(.horizontal, 8)
        .padding(.top, 12)
        .padding(.bottom, 4)
    }

    private func serverRow(_ server: VPNServer) -> some View {
        let isSelected = server.id == vm.selectedServer?.id
        let isAuto = server.id == "auto"

        return Button(action: {
            if isAuto { vm.selectFastest() }
            else { vm.selectedServer = server }
            dismiss()
        }) {
            HStack(spacing: 12) {
                if !isAuto {
                    Text(server.countryFlag)
                        .font(.title3)
                } else {
                    Image(systemName: "bolt.fill")
                        .font(.system(size: 16))
                        .foregroundColor(.accentColor)
                        .frame(width: 32, height: 32)
                        .background(Color.accentColor.opacity(0.1), in: RoundedRectangle(cornerRadius: 8))
                }

                VStack(alignment: .leading, spacing: 2) {
                    Text(server.city)
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundColor(.primary)
                    Text(isAuto ? "Auto-selected based on lowest latency" : "\(server.latencyMs)ms · \(server.loadPercent)% · \(server.name)")
                        .font(.system(size: 10))
                        .foregroundColor(.secondary)
                }

                Spacer()

                if server.isPremium {
                    Text("PREMIUM")
                        .font(.system(size: 7, weight: .bold))
                        .foregroundColor(.white)
                        .padding(.horizontal, 4)
                        .padding(.vertical, 2)
                        .background(Color.orange, in: RoundedRectangle(cornerRadius: 3))
                }

                if isSelected {
                    Image(systemName: "checkmark.circle.fill")
                        .font(.system(size: 14))
                        .foregroundColor(.accentColor)
                }
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 10)
            .background(isSelected ? Color.accentColor.opacity(0.08) : Color.clear, in: RoundedRectangle(cornerRadius: 8))
            .contentShape(Rectangle())
        }
        .buttonStyle(.plain)
    }
}
