import SwiftUI

struct OnboardingView: View {
    @EnvironmentObject private var vm: MainViewModel
    @State private var currentPage = 0

    let pages: [(icon: String, title: String, text: String)] = [
        ("shield.checkered", "Secure Connection",
         "Your traffic is encrypted and routed through our global network of servers."),
        ("globe", "Global Network",
         "Connect to 50+ servers in 30 countries. Always find the fastest route."),
        ("bolt.shield", "Anti-DPI Technology",
         "Advanced obfuscation defeats deep packet inspection. TLS camouflage, packet fragmentation, and more."),
        ("sparkles", "Ready to Go",
         "You're all set. Connect with one tap and browse freely."),
    ]

    var body: some View {
        VStack(spacing: 0) {
            Spacer()

            // Page content
            TabView(selection: $currentPage) {
                ForEach(0..<pages.count, id: \.self) { index in
                    VStack(spacing: 20) {
                        Image(systemName: pages[index].icon)
                            .font(.system(size: 52))
                            .foregroundColor(.accentColor)
                        Text(pages[index].title)
                            .font(.system(size: 22, weight: .bold))
                        Text(pages[index].text)
                            .font(.system(size: 13))
                            .foregroundColor(.secondary)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal, 40)
                            .fixedSize(horizontal: false, vertical: true)
                    }
                    .tag(index)
                }
            }
            .tabViewStyle(.automatic)

            Spacer()

            // Page indicator
            HStack(spacing: 8) {
                ForEach(0..<pages.count, id: \.self) { index in
                    Circle()
                        .fill(currentPage == index ? Color.accentColor : Color.secondary.opacity(0.3))
                        .frame(width: 8, height: 8)
                }
            }
            .padding(.bottom, 24)

            // Buttons
            HStack(spacing: 12) {
                if currentPage < pages.count - 1 {
                    Button("Skip") {
                        vm.completeOnboarding()
                    }
                    .buttonStyle(.plain)
                    .foregroundColor(.secondary)
                }

                Spacer()

                Button(currentPage < pages.count - 1 ? "Next" : "Get Started") {
                    if currentPage < pages.count - 1 {
                        withAnimation { currentPage += 1 }
                    } else {
                        vm.completeOnboarding()
                    }
                }
                .buttonStyle(.borderedProminent)
            }
            .padding(.horizontal, 40)
            .padding(.bottom, 40)
        }
        .frame(width: 380, height: 520)
    }
}
