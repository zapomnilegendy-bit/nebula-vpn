'use client'

import { ArrowDown, ArrowUp, Shield, Globe, Zap } from 'lucide-react'

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-[var(--text-primary)]">Dashboard</h1>
        <p className="text-sm text-[var(--text-secondary)] mt-1">Welcome back, focus — your connection is secured.</p>
      </div>

      {/* Connection Card */}
      <div className="card p-8">
        <div className="flex items-center gap-6">
          {/* Status ring */}
          <div className="relative w-28 h-28 flex-shrink-0">
            <svg className="w-full h-full -rotate-90" viewBox="0 0 120 120">
              <circle cx="60" cy="60" r="52" fill="none" stroke="currentColor" strokeWidth="4" className="text-nebula-200 dark:text-nebula-700" />
              <circle cx="60" cy="60" r="52" fill="none" stroke="#34C759" strokeWidth="4" strokeDasharray="320" strokeDashoffset="80" strokeLinecap="round" />
            </svg>
            <div className="absolute inset-0 flex items-center justify-center">
              <Shield className="w-10 h-10 text-[#34C759]" />
            </div>
          </div>

          {/* Status info */}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <span className="w-2.5 h-2.5 rounded-full bg-[#34C759]" />
              <span className="text-sm font-semibold text-[#34C759]">Connected</span>
            </div>
            <h2 className="text-xl font-bold text-[var(--text-primary)]">EU-Frankfurt-01</h2>
            <p className="text-sm text-[var(--text-secondary)] mt-1">WireGuard · 12ms latency · 0h 42m 15s</p>

            {/* Data transfer */}
            <div className="flex gap-6 mt-4">
              <div className="flex items-center gap-2 text-sm text-[var(--text-secondary)]">
                <ArrowDown className="w-3.5 h-3.5" />
                <span className="font-mono">1.2 GB</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-[var(--text-secondary)]">
                <ArrowUp className="w-3.5 h-3.5" />
                <span className="font-mono">342 MB</span>
              </div>
            </div>
          </div>

          {/* Disconnect button */}
          <button className="btn-primary bg-red-500 hover:bg-red-600 text-white">
            Disconnect
          </button>
        </div>
      </div>

      {/* Quick actions */}
      <div className="grid grid-cols-3 gap-4">
        <div className="card p-5 hover:shadow-apple-lg transition-shadow cursor-pointer">
          <Globe className="w-5 h-5 text-[var(--text-secondary)] mb-3" />
          <h3 className="text-sm font-semibold text-[var(--text-primary)]">Fastest Server</h3>
          <p className="text-xs text-[var(--text-secondary)] mt-1">Auto-selected: Frankfurt</p>
        </div>
        <div className="card p-5 hover:shadow-apple-lg transition-shadow cursor-pointer">
          <Shield className="w-5 h-5 text-[var(--text-secondary)] mb-3" />
          <h3 className="text-sm font-semibold text-[var(--text-primary)]">Kill Switch</h3>
          <p className="text-xs text-[var(--text-secondary)] mt-1">Enabled — traffic secured</p>
        </div>
        <div className="card p-5 hover:shadow-apple-lg transition-shadow cursor-pointer">
          <Zap className="w-5 h-5 text-[var(--text-secondary)] mb-3" />
          <h3 className="text-sm font-semibold text-[var(--text-primary)]">Protocol</h3>
          <p className="text-xs text-[var(--text-secondary)] mt-1">WireGuard · Optimal</p>
        </div>
      </div>

      {/* Sessions today */}
      <div>
        <h2 className="text-sm font-semibold text-[var(--text-secondary)] mb-3 uppercase tracking-wide">Recent Sessions</h2>
        <div className="card divide-y divide-black/5 dark:divide-white/5">
          {[
            { server: 'EU-Frankfurt-01', location: 'Frankfurt, DE', protocol: 'WireGuard', duration: '0h 42m', data: '1.5 GB', status: 'active' as const },
            { server: 'US-NewYork-01', location: 'New York, US', protocol: 'VLESS+Reality', duration: '1h 12m', data: '3.2 GB', status: 'ended' as const },
            { server: 'JP-Tokyo-01', location: 'Tokyo, JP', protocol: 'Hysteria2', duration: '0h 28m', data: '890 MB', status: 'ended' as const },
          ].map((session, i) => (
            <div key={i} className="flex items-center justify-between px-5 py-4">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-nebula-100 dark:bg-nebula-800 flex items-center justify-center text-sm">
                  {session.location.split(',')[0][0]}
                </div>
                <div>
                  <p className="text-sm font-medium text-[var(--text-primary)]">{session.server}</p>
                  <p className="text-xs text-[var(--text-secondary)]">{session.location}</p>
                </div>
              </div>
              <div className="flex items-center gap-6 text-sm">
                <span className="text-[var(--text-secondary)] font-mono">{session.duration}</span>
                <span className="text-[var(--text-secondary)] font-mono">{session.data}</span>
                <span className={session.status === 'active' ? 'badge-green' : 'badge-orange'}>
                  {session.status === 'active' ? 'Active' : 'Ended'}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
