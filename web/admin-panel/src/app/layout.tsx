export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="dark">
      <head><title>NebulaVPN Admin</title></head>
      <body className="min-h-screen bg-[var(--bg-secondary)]">{children}</body>
    </html>
  )
}
