import SwiftUI

// MARK: - Toast Model
struct Toast: Equatable {
    let id = UUID()
    let message: String
    let type: ToastType

    enum ToastType: Equatable {
        case success, error, info, warning
    }

    var icon: String {
        switch type {
        case .success: return "checkmark.circle.fill"
        case .error:   return "xmark.circle.fill"
        case .info:    return "info.circle.fill"
        case .warning: return "exclamationmark.triangle.fill"
        }
    }

    var color: Color {
        switch type {
        case .success: return Color(red: 0.22, green: 0.78, blue: 0.35)
        case .error:   return .red
        case .info:    return .blue
        case .warning: return .orange
        }
    }
}

// MARK: - Toast View
struct ToastView: View {
    let toast: Toast

    var body: some View {
        HStack(spacing: 10) {
            Image(systemName: toast.icon)
                .font(.system(size: 16, weight: .medium))
                .foregroundColor(toast.color)
            Text(toast.message)
                .font(.system(size: 14, weight: .medium))
                .foregroundColor(.primary)
            Spacer()
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 12)
        .background(.ultraThinMaterial, in: RoundedRectangle(cornerRadius: 12))
        .overlay(RoundedRectangle(cornerRadius: 12).stroke(Color.primary.opacity(0.06), lineWidth: 1))
        .shadow(color: .black.opacity(0.08), radius: 12, y: 4)
        .padding(.horizontal, 16)
        .transition(.move(edge: .top).combined(with: .opacity))
    }
}

// MARK: - Toast Modifier
struct ToastModifier: ViewModifier {
    @Binding var toast: Toast?
    @State private var workItem: DispatchWorkItem?

    func body(content: Content) -> some View {
        content
            .overlay(alignment: .top) {
                if let toast = toast {
                    ToastView(toast: toast)
                        .padding(.top, 8)
                        .onAppear { dismissAfterDelay() }
                        .onTapGesture { dismiss() }
                }
            }
            .animation(.spring(response: 0.35), value: toast?.id)
    }

    private func dismissAfterDelay() {
        workItem?.cancel()
        let item = DispatchWorkItem { dismiss() }
        workItem = item
        DispatchQueue.main.asyncAfter(deadline: .now() + 3, execute: item)
    }

    private func dismiss() {
        withAnimation { toast = nil }
        workItem?.cancel()
        workItem = nil
    }
}

// MARK: - View Extension
extension View {
    func toast(toast: Binding<Toast?>) -> some View {
        modifier(ToastModifier(toast: toast))
    }
}

// MARK: - Confirmable Alert
struct ConfirmAlert: Identifiable {
    let id = UUID()
    let title: String
    let message: String
    let confirmLabel: String
    let cancelLabel: String
    let destructive: Bool
    let action: () -> Void
}

extension View {
    func confirmAlert(alert: Binding<ConfirmAlert?>) -> some View {
        alert(
            alert.wrappedValue?.title ?? "",
            isPresented: .init(get: { alert.wrappedValue != nil }, set: { if !$0 { alert.wrappedValue = nil } }),
            actions: {
                Button(alert.wrappedValue?.cancelLabel ?? "Cancel", role: .cancel) { alert.wrappedValue = nil }
                Button(alert.wrappedValue?.confirmLabel ?? "OK", role: alert.wrappedValue?.destructive == true ? .destructive : nil) {
                    alert.wrappedValue?.action()
                    alert.wrappedValue = nil
                }
            },
            message: { Text(alert.wrappedValue?.message ?? "") }
        )
    }
}
