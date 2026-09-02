'use client'

import AppShell from '@/components/AppShell'
import { Globe, Signal, Wifi } from 'lucide-react'

const servers = [
  { id: 'eu-fra-01', name: 'EU-Frankfurt-01', country: '🇩🇪', city: 'Frankfurt', latency: 12, load: 42, protocol: 'WireGuard', premium: false, healthy: true },
  { id: 'eu-ams-01', name: 'EU-Amsterdam-01', country: '🇳🇱', city: 'Amsterdam', latency: 18, load: 58, protocol: 'WireGuard', premium: false, healthy: true },
  { id: 'gb-lon-01', name: 'GB-London-01', country: '🇬🇧', city: 'London', latency: 24, load: 55, protocol: 'VLESS+Reality', premium: false, healthy: true },
  { id: 'us-nyc-01', name: 'US-NewYork-01', country: '🇺🇸', city: 'New York', latency: 89, load: 71, protocol: 'WireGuard', premium: true, healthy: true },
  { id: 'us-sfo-01', name: 'US-SanFrancisco-01', country: '🇺🇸', city: 'San Francisco', latency: 105, load: 34, protocol: 'Hysteria2', premium: false, healthy: true },
  { id: 'us-mia-01', name: 'US-Miami-01', country: '🇺🇸', city: 'Miami', latency: 130, load: 29, protocol: 'WireGuard', premium: false, healthy: true },
  { id: 'jp-tokyo-01', name: 'JP-Tokyo-01', country: '🇯🇵', city: 'Tokyo', latency: 185, load: 23, protocol: 'TUIC', premium: true, healthy: true },
  { id: 'sg-sin-01', name: 'SG-Singapore-01', country: '🇸🇬', city: 'Singapore', latency: 210, load: 18, protocol: 'WireGuard', premium: false, healthy: true },
  { id: 'au-syd-01', name: 'AU-Sydney-01', country: '🇦🇺', city: 'Sydney', latency: 280, load: 12, protocol: 'Hysteria2', premium: true, healthy: true },
  { id: 'br-gru-01', name: 'BR-SaoPaulo-01', country: '🇧🇷', city: 'São Paulo', latency: 175, load: 9, protocol: 'VLESS+Reality', premium: false, healthy: true },
]

const regions = ['All', 'Europe', 'Americas', 'Asia Pacific']

export default function ServersPage() {
  return (
    <AppShell>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-[var(--text-primary)]">Servers</h1>
          <p className="text-sm text-[var(--text-secondary)] mt-1">Select a server to connect through.</p>
        </div>

        {/* Region filter */}
        <div className="flex gap-2">
          {regions.map(r => (
            <button key={r} className={`px-4 py-2 rounded-apple text-sm font-medium transition-all ${
              r === 'All' ? 'bg-black dark:bg-white text-white dark:text-black' : 'bg-nebula-100 dark:bg-nebula-800 text-[var(--text-secondary)] hover:text-[var(--text-primary)]'
            }`}>
              {r}
            </button>
          ))}
        </div>

        {/* Server list */}
        <div className="space-y-2">
          {servers.map(server => (
            <div key={server.id} className="card p-4 flex items-center justify-between hover:shadow-apple-lg transition-shadow cursor-pointer group">
              <div className="flex items-center gap-4">
                <span className="text-2xl">{server.country}</span>
                <div>
                  <div className="flex items-center gap-2">
                    <h3 className="text-sm font-semibold text-[var(--text-primary)]">{server.city}</h3>
                    {server.premium && <span className="px-1.5 py-0.5 rounded text-[10px] font-bold bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400">PREMIUM</span>}
                  </div>
                  <p className="text-xs text-[var(--text-secondary)] font-mono mt-0.5">{server.name}</p>
                </div>
              </div>

              <div className="flex items-center gap-6">
                <div className="text-right">
                  <p className="text-sm font-mono text-[var(--text-primary)]">{server.latency}ms</p>
                  <p className="text-xs text-[var(--text-secondary)]">{server.protocol}</p>
                </div>
                <div className="w-16">
                  <div className="h-1.5 rounded-full bg-nebula-200 dark:bg-nebula-700 overflow-hidden">
                    <div className={`h-full rounded-full transition-all ${
                      server.load < 40 ? 'bg-[#34C759]' : server.load < 70 ? 'bg-[#FF9500]' : 'bg-[#FF3B30]'
                    }`} style={{ width: `${server.load}%` }} />
                  </div>
                </div>
                <button className="btn-primary opacity-0 group-hover:opacity-100 transition-opacity text-xs py-2 px-4">
                  Connect
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </AppShell>
  )
}
