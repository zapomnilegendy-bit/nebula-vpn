import SwiftUI

struct ConnectionRingView: View {
    let status: ConnectionStatus
    let progress: CGFloat

    @State private var rotation: Double = 0
    @State private var breathe: CGFloat = 1

    var body: some View {
        ZStack {
            // Outer glow
            if status == .connected {
                Circle()
                    .fill(Color(red: 0.22, green: 0.78, blue: 0.35).opacity(0.1))
                    .scaleEffect(breathe)
                    .animation(.easeInOut(duration: 2.5).repeatForever(autoreverses: true), value: breathe)
                    .onAppear { breathe = 1.1 }
            }

            // Background ring
            Circle()
                .stroke(lineWidth: 5)
                .foregroundColor(.primary.opacity(0.06))

            // Active ring
            Circle()
                .trim(from: 0, to: progress)
                .stroke(
                    AngularGradient(
                        gradient: Gradient(colors: status == .connected
                            ? [Color(red: 0.22, green: 0.78, blue: 0.35), Color(red: 0.3, green: 0.5, blue: 1.0)]
                            : [Color.primary.opacity(0.15), Color.primary.opacity(0.05)]),
                        center: .center
                    ),
                    style: StrokeStyle(lineWidth: 5, lineCap: .round)
                )
                .rotationEffect(.degrees(-90 + rotation))
                .animation(.linear(duration: 4).repeatForever(autoreverses: false), value: rotation)
                .onAppear { rotation = 360 }

            // Center icon
            Image(systemName: status.icon)
                .font(.system(size: 44, weight: .regular))
                .foregroundColor(status.color)
        }
    }
}
