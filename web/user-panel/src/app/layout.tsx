export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className="dark">
      <head>
        <title>NebulaVPN</title>
        <meta name="description" content="Secure, fast, and private VPN service" />
        <link rel="icon" href="/favicon.ico" />
      </head>
      <body className="min-h-screen">{children}</body>
    </html>
  )
}
