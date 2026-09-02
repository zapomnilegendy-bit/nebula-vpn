import SwiftUI

struct AboutView: View {
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        VStack(spacing: 20) {
            Image(systemName: "shield.checkered")
                .font(.system(size: 48))
                .foregroundColor(.accentColor)

            Text("NebulaVPN")
                .font(.system(size: 20, weight: .bold))

            Text("Version 1.0.0 (build 1)")
                .font(.system(size: 12))
                .foregroundColor(.secondary)

            Text("© 2026 NebulaVPN. All rights reserved.\nSecure, private, unrestricted internet access.")
                .font(.system(size: 11))
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)

            Button("Done") { dismiss() }
                .buttonStyle(.borderedProminent)
                .controlSize(.small)
                .padding(.top, 8)
        }
        .padding(32)
        .frame(width: 260)
    }
}
