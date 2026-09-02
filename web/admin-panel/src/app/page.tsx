'use client'

import AdminShell from '@/components/AdminShell'
import { Users, Server, CreditCard, Activity, TrendingUp, Zap, Globe, Shield } from 'lucide-react'

const stats = [
  { icon: Users, label: 'Total Users', value: '12,847', change: '+12.3%', color: 'text-accent' },
  { icon: Server, label: 'Active Nodes', value: '47', change: '+3', color: 'text-[#34C759]' },
  { icon: CreditCard, label: 'MRR', value: '$84,293', change: '+8.7%', color: 'text-[#FF9500]' },
  { icon: Activity, label: 'Avg. Latency', value: '42ms', change: '-5ms', color: 'text-[#34C759]' },
]

const recentUsers = [
  { email: 'alice@example.com', plan: 'Pro', status: 'active', created: '2m ago', data: '1.2 GB' },
  { email: 'bob@corp.com', plan: 'Ultimate', status: 'active', created: '15m ago', data: '4.7 GB' },
  { email: 'charlie@dev.io', plan: 'Free', status: 'active', created: '1h ago', data: '120 MB' },
  { email: 'diana@design.co', plan: 'Pro', status: 'suspended', created: '2h ago', data: '890 MB' },
]

export default function AdminDashboard() {
  return (
    <AdminShell>
      <div className="space-y-8">
        <div>
          <h1 className="text-xl font-bold text-[var(--text-primary)]">Dashboard</h1>
          <p className="text-sm text-[var(--text-secondary)] mt-1">System overview and key metrics.</p>
        </div>

        {/* Stats grid */}
        <div className="grid grid-cols-4 gap-4">
          {stats.map((s, i) => (
            <div key={i} className="card p-5">
              <div className="flex items-center justify-between mb-3">
                <s.icon className={`w-5 h-5 ${s.color}`} />
                <span className={`text-xs font-semibold ${s.change.startsWith('+') ? 'text-[#34C759]' : 'text-[#FF3B30]'}`}>
                  {s.change}
                </span>
              </div>
              <p className="text-2xl font-bold text-[var(--text-primary)]">{s.value}</p>
              <p className="text-xs text-[var(--text-secondary)] mt-1">{s.label}</p>
            </div>
          ))}
        </div>

        {/* Charts */}
        <div className="grid grid-cols-2 gap-4">
          <div className="card p-5">
            <h3 className="text-sm font-semibold text-[var(--text-primary)] mb-4">Connections (24h)</h3>
            <div className="h-48 flex items-end gap-1">
              {[30, 45, 25, 60, 80, 55, 70, 90, 65, 40, 75, 85, 95, 100, 80, 60, 50, 70, 85, 90, 75, 55, 40, 30].map((h, i) => (
                <div key={i} className="flex-1 rounded-t-sm bg-accent/20 dark:bg-accent/10" style={{ height: `${h}%` }} />
              ))}
            </div>
          </div>
          <div className="card p-5">
            <h3 className="text-sm font-semibold text-[var(--text-primary)] mb-4">Bandwidth Usage</h3>
            <div className="space-y-3">
              {[
                { label: 'WireGuard', value: 65, color: 'bg-[#34C759]' },
                { label: 'VLESS+Reality', value: 20, color: 'bg-accent' },
                { label: 'Hysteria2', value: 10, color: 'bg-[#FF9500]' },
                { label: 'TUIC', value: 5, color: 'bg-[#FF3B30]' },
              ].map((p, i) => (
                <div key={i}>
                  <div className="flex justify-between text-xs mb-1">
                    <span className="text-[var(--text-primary)]">{p.label}</span>
                    <span className="text-[var(--text-secondary)]">{p.value}%</span>
                  </div>
                  <div className="h-2 rounded-full bg-nebula-100 dark:bg-nebula-800 overflow-hidden">
                    <div className={`h-full rounded-full ${p.color}`} style={{ width: `${p.value}%` }} />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Recent users */}
        <div className="card">
          <div className="px-5 py-4 border-b border-black/5 dark:border-white/5 flex items-center justify-between">
            <h3 className="text-sm font-semibold text-[var(--text-primary)]">Recent Users</h3>
            <button className="btn-secondary text-xs">View All</button>
          </div>
          <div className="divide-y divide-black/5 dark:divide-white/5">
            {recentUsers.map((u, i) => (
              <div key={i} className="flex items-center justify-between px-5 py-3.5">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-nebula-100 dark:bg-nebula-800 flex items-center justify-center text-xs font-semibold">
                    {u.email[0].toUpperCase()}
                  </div>
                  <div>
                    <p className="text-sm font-medium text-[var(--text-primary)]">{u.email}</p>
                    <p className="text-xs text-[var(--text-secondary)]">Joined {u.created} · {u.data}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <span className={`px-2 py-0.5 rounded text-[10px] font-semibold ${
                    u.plan === 'Ultimate' ? 'bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400' :
                    u.plan === 'Pro' ? 'bg-accent/10 text-accent' : 'bg-nebula-100 dark:bg-nebula-800 text-[var(--text-secondary)]'
                  }`}>{u.plan}</span>
                  <span className={`px-2 py-0.5 rounded text-[10px] font-semibold ${
                    u.status === 'active' ? 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400' :
                    'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400'
                  }`}>{u.status}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </AdminShell>
  )
}
