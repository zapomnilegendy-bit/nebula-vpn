import SwiftUI

struct OnboardingView: View {
    @EnvironmentObject private var appState: AppState
    @State private var page = 0

    let pages: [(icon: String, title: String, desc: String)] = [
        ("shield.checkered", "Privacy First",
         "NebulaVPN encrypts all your traffic with military-grade WireGuard protocol. No logs, no tracking."),
        ("globe", "Global Network",
         "Access content from 50+ locations worldwide. Your real IP stays hidden everywhere you go."),
        ("bolt.shield", "Anti-DPI Technology",
         "Advanced obfuscation bypasses deep packet inspection. Works in the most restrictive networks."),
        ("qrcode.viewfinder", "One-Click Connect",
         "Import config via QR code or select from our optimized server list. Connect in one click."),
    ]

    var body: some View {
        VStack(spacing: 0) {
            TabView(selection: $page) {
                ForEach(Array(pages.enumerated()), id: \.offset) { i, p in
                    VStack(spacing: 24) {
                        Spacer()
                        Image(systemName: p.icon)
                            .font(.system(size: 72))
                            .foregroundColor(.accentColor)
                        VStack(spacing: 8) {
                            Text(p.title)
                                .font(.system(size: 24, weight: .bold))
                            Text(p.desc)
                                .font(.system(size: 15))
                                .foregroundColor(.secondary)
                                .multilineTextAlignment(.center)
                                .lineSpacing(4)
                                .padding(.horizontal, 32)
                        }
                        Spacer()
                    }
                    .tag(i)
                }
            }
            .tabViewStyle(.page(indexDisplayMode: .always))

            HStack(spacing: 12) {
                Button("Skip") {
                    appState.completeOnboarding()
                }
                .font(.system(size: 15))
                .foregroundColor(.secondary)

                Spacer()

                Button(page < pages.count - 1 ? "Continue" : "Get Started") {
                    if page < pages.count - 1 {
                        withAnimation { page += 1 }
                    } else {
                        appState.completeOnboarding()
                    }
                }
                .buttonStyle(.borderedProminent)
                .controlSize(.large)
            }
            .padding(24)
        }
    }
}
