// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "NebulaVPN",
    platforms: [
        .macOS(.v14)
    ],
    products: [
        .executable(name: "NebulaVPN", targets: ["NebulaVPN"])
    ],
    targets: [
        .executableTarget(
            name: "NebulaVPN",
            path: "Sources/NebulaVPN"
        )
    ]
)
