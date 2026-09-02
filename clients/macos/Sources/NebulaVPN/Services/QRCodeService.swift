import Foundation
import CoreImage
import AppKit

final class QRCodeService {
    func generateQRCode(from configString: String) -> NSImage? {
        guard let filter = CIFilter(name: "CIQRCodeGenerator") else { return nil }

        let data = Data(configString.utf8)
        filter.setValue(data, forKey: "inputMessage")
        filter.setValue("M", forKey: "inputCorrectionLevel")

        guard let ciImage = filter.outputImage else { return nil }

        let transform = CGAffineTransform(scaleX: 10, y: 10)
        let scaled = ciImage.transformed(by: transform)

        let rep = NSCIImageRep(ciImage: scaled)
        let nsImage = NSImage(size: rep.size)
        nsImage.addRepresentation(rep)

        return nsImage
    }

    func readQRCode(from image: NSImage) -> String? {
        guard let cgImage = image.cgImage(forProposedRect: nil, context: nil, hints: nil) else { return nil }

        let ciImage = CIImage(cgImage: cgImage)
        let detector = CIDetector(ofType: CIDetectorTypeQRCode,
                                  context: nil,
                                  options: [CIDetectorAccuracy: CIDetectorAccuracyHigh])

        guard let features = detector?.features(in: ciImage) as? [CIQRCodeFeature],
              let first = features.first else { return nil }

        return first.messageString
    }

    func parseConfig(_ configString: String) -> (host: String, port: Int, publicKey: String)? {
        let lines = configString.components(separatedBy: .newlines)
        var host = ""
        var port = 51820
        var publicKey = ""

        for line in lines {
            let trimmed = line.trimmingCharacters(in: .whitespaces)
            if trimmed.hasPrefix("Endpoint =") {
                let endpoint = trimmed.replacingOccurrences(of: "Endpoint = ", with: "")
                let parts = endpoint.components(separatedBy: ":")
                host = parts.first ?? ""
                if parts.count > 1 { port = Int(parts[1]) ?? 51820 }
            }
            if trimmed.hasPrefix("PublicKey =") {
                publicKey = trimmed.replacingOccurrences(of: "PublicKey = ", with: "")
            }
        }

        guard !host.isEmpty, !publicKey.isEmpty else { return nil }
        return (host, port, publicKey)
    }
}
