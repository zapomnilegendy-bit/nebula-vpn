import SwiftUI

struct ContentView: View {
    @EnvironmentObject private var appState: AppState
    @State private var pulsing = false

    var body: some View {
        ZStack {
            if appState.isOnboarding {
                OnboardingView()
            } else {
                mainContent
            }
        }
    }

    private var mainContent: some View {
        ZStack {
            backgroundGradient
                .ignoresSafeArea()

            VStack(spacing: 0) {
                // Top bar with settings and about
                HStack {
                    Spacer()
                    Button(action: { appState.showingAbout = true }) {
                        Image(systemName: "info.circle")
                            .font(.system(size: 15, weight: .regular))
                            .foregroundColor(.primary.opacity(0.6))
                            .frame(width: 36, height: 36)
                            .background(.ultraThinMaterial, in: Circle())
                    }
                    Button(action: { appState.showingSettings = true }) {
                        Image(systemName: "gearshape")
                            .font(.system(size: 17, weight: .regular))
                            .foregroundColor(.primary.opacity(0.6))
                            .frame(width: 36, height: 36)
                            .background(.ultraThinMaterial, in: Circle())
                    }
                }
                .padding(.horizontal, 20)
                .padding(.top, 8)

                Spacer()

                // Status shield
                ZStack {
                    // Outer glow
                    if appState.isConnected {
                        Circle()
                            .fill(Color(red: 0.22, green: 0.78, blue: 0.35).opacity(0.15))
                            .frame(width: 260, height: 260)
                            .scaleEffect(pulsing ? 1.1 : 0.95)
                            .animation(.easeInOut(duration: 2).repeatForever(autoreverses: true), value: pulsing)
                    }

                    // Main ring
                    Circle()
                        .stroke(lineWidth: 3)
                        .foregroundColor(.primary.opacity(0.08))
                        .frame(width: 200, height: 200)

                    // Active ring
                    Circle()
                        .trim(from: 0, to: appState.isConnected ? 0.85 : 0.15)
                        .stroke(
                            AngularGradient(
                                gradient: Gradient(colors: appState.isConnected
                                    ? [Color(red: 0.22, green: 0.78, blue: 0.35), Color(red: 0.3, green: 0.5, blue: 1.0)]
                                    : [Color.primary.opacity(0.15), Color.primary.opacity(0.05)]
                                ),
                                center: .center
                            ),
                            style: StrokeStyle(lineWidth: 3, lineCap: .round)
                        )
                        .frame(width: 200, height: 200)
                        .rotationEffect(.degrees(-90))
                        .animation(.spring(response: 0.8), value: appState.isConnected)

                    // Shield icon
                    Image(systemName: appState.statusIcon)
                        .font(.system(size: 56, weight: .regular))
                        .foregroundColor(appState.isConnected ? Color(red: 0.22, green: 0.78, blue: 0.35) : .primary.opacity(0.25))
                        .symbolEffect(.bounce, value: appState.isConnected)
                }

                Spacer().frame(height: 40)

                // Status text
                Text(appState.statusText)
                    .font(.system(size: 20, weight: .semibold))
                    .foregroundColor(.primary)

                if appState.isConnected {
                    Text(appState.formattedTime)
                        .font(.system(size: 15, weight: .regular, design: .monospaced))
                        .foregroundColor(.secondary)
                        .padding(.top, 4)
                }

                Spacer().frame(height: 32)

                // Data usage
                if appState.isConnected {
                    HStack(spacing: 32) {
                        DataBadge(icon: "arrow.down", value: appState.dataTransferred.formattedDownload)
                        DataBadge(icon: "arrow.up", value: appState.dataTransferred.formattedUpload)
                    }
                    .transition(.move(edge: .bottom).combined(with: .opacity))
                }

                Spacer()

                // Server selector
                VStack(spacing: 12) {
                    Button(action: { appState.showingServerList = true }) {
                        HStack(spacing: 12) {
                            if let server = appState.selectedServer {
                                Text(server.countryFlag)
                                    .font(.title2)
                                VStack(alignment: .leading, spacing: 2) {
                                    Text(server.city)
                                        .font(.system(size: 15, weight: .semibold))
                                        .foregroundColor(.primary)
                                    Text("\(server.latencyMs)ms · \(server.loadPercent)% load")
                                        .font(.system(size: 12))
                                        .foregroundColor(.secondary)
                                }
                            } else {
                                Image(systemName: "magnifyingglass")
                                    .font(.system(size: 15, weight: .medium))
                                    .foregroundColor(.secondary)
                                Text("Select server")
                                    .font(.system(size: 15, weight: .medium))
                                    .foregroundColor(.secondary)
                            }
                            Spacer()
                            Image(systemName: "chevron.right")
                                .font(.system(size: 12, weight: .semibold))
                                .foregroundColor(.secondary)
                        }
                        .padding(.horizontal, 16)
                        .padding(.vertical, 14)
                        .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 14))
                    }
                    .buttonStyle(.plain)
                    .padding(.horizontal, 20)

                    // Connect button
                    Button(action: {
                        let impact = UIImpactFeedbackGenerator(style: .medium)
                        impact.impactOccurred()
                        if appState.isConnected {
                            appState.disconnect()
                        } else {
                            appState.connect()
                        }
                    }) {
                        ZStack {
                            RoundedRectangle(cornerRadius: 14)
                                .fill(appState.isConnected ? Color.red : Color.primary)
                                .frame(height: 50)

                            Text(appState.isConnected ? "Disconnect" : "Connect")
                                .font(.system(size: 16, weight: .semibold))
                                .foregroundColor(appState.isConnected ? .white : Color(.systemBackground))
                        }
                    }
                    .buttonStyle(.plain)
                    .padding(.horizontal, 20)
                }

                Spacer().frame(height: 24)
            }
        }
        .onAppear {
            pulsing = true
        }
        .sheet(isPresented: $appState.showingServerList) {
            ServerListView()
        }
        .sheet(isPresented: $appState.showingSettings) {
            SettingsView()
        }
        .sheet(isPresented: $appState.showingAbout) {
            AboutView()
        }
    }

    private var backgroundGradient: some View {
        LinearGradient(
            gradient: Gradient(colors: [
                Color(.systemBackground),
                Color(.systemGray6).opacity(0.8)
            ]),
            startPoint: .top,
            endPoint: .bottom
        )
    }
}

struct DataBadge: View {
    let icon: String
    let value: String

    var body: some View {
        HStack(spacing: 6) {
            Image(systemName: icon)
                .font(.system(size: 11, weight: .semibold))
                .foregroundColor(.secondary)
            Text(value)
                .font(.system(size: 13, weight: .medium, design: .monospaced))
                .foregroundColor(.secondary)
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 8)
        .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 10))
    }
}
