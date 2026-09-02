import SwiftUI

struct AboutView: View {
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        NavigationStack {
            VStack(spacing: 24) {
                Image(systemName: "shield.checkered")
                    .font(.system(size: 64))
                    .foregroundColor(.accentColor)

                Text("NebulaVPN")
                    .font(.system(size: 24, weight: .bold))

                Text("Version 1.0.0 (build 1)")
                    .font(.system(size: 14))
                    .foregroundColor(.secondary)

                Text("© 2026 NebulaVPN. All rights reserved.\nSecure, private, unrestricted internet access.")
                    .font(.system(size: 13))
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
            }
            .padding(32)
            .toolbar {
                ToolbarItem(placement: .confirmationAction) {
                    Button("Done") { dismiss() }
                        .font(.system(size: 15, weight: .semibold))
                }
            }
        }
    }
}
