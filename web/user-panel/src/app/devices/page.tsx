'use client'

import AppShell from '@/components/AppShell'
import { Monitor, Laptop, Smartphone, Tablet, Plus, MoreHorizontal } from 'lucide-react'

const devices = [
  { name: "Focus's MacBook Pro", platform: 'macos', icon: Laptop, lastActive: '2 min ago', ip: '10.0.0.2', data: '1.2 GB', active: true },
  { name: "Focus's iPhone 15", platform: 'ios', icon: Smartphone, lastActive: '15 min ago', ip: '10.0.0.3', data: '890 MB', active: true },
  { name: "Focus's iPad Air", platform: 'ios', icon: Tablet, lastActive: '2h ago', ip: '10.0.0.4', data: '340 MB', active: false },
]

const platformColors: Record<string, string> = {
  macos: 'bg-[#007AFF]',
  ios: 'bg-[#34C759]',
  windows: 'bg-[#FF9500]',
  android: 'bg-[#34C759]',
}

export default function DevicesPage() {
  return (
    <AppShell>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-[var(--text-primary)]">Devices</h1>
            <p className="text-sm text-[var(--text-secondary)] mt-1">Manage your connected devices. 2 of 5 active.</p>
          </div>
          <button className="btn-primary flex items-center gap-2">
            <Plus className="w-4 h-4" />
            Add Device
          </button>
        </div>

        {/* Device list */}
        <div className="space-y-3">
          {devices.map((device, i) => {
            const Icon = device.icon
            return (
              <div key={i} className="card p-5">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className={`w-10 h-10 rounded-apple ${platformColors[device.platform]} flex items-center justify-center`}>
                      <Icon className="w-5 h-5 text-white" />
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <h3 className="text-sm font-semibold text-[var(--text-primary)]">{device.name}</h3>
                        {device.active && <span className="w-2 h-2 rounded-full bg-[#34C759]" />}
                      </div>
                      <p className="text-xs text-[var(--text-secondary)] mt-0.5">
                        {device.ip} · {device.data} used · {device.lastActive}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <button className="btn-secondary text-xs py-2 px-3">
                      {device.active ? 'Disconnect' : 'Connect'}
                    </button>
                    <button className="p-2 rounded-apple hover:bg-nebula-100 dark:hover:bg-nebula-800 transition-all">
                      <MoreHorizontal className="w-4 h-4 text-[var(--text-secondary)]" />
                    </button>
                  </div>
                </div>
              </div>
            )
          })}
        </div>

        {/* How to add device */}
        <div className="card p-6">
          <h3 className="text-sm font-semibold text-[var(--text-primary)] mb-3">Add a New Device</h3>
          <div className="grid grid-cols-3 gap-4">
            {[
              { platform: 'macOS', icon: Laptop, desc: 'Download for Mac' },
              { platform: 'Windows', icon: Monitor, desc: 'Download for Windows' },
              { platform: 'iOS', icon: Smartphone, desc: 'Download on App Store' },
              { platform: 'Android', icon: Smartphone, desc: 'Get it on Google Play' },
            ].map((item, i) => (
              <div key={i} className="flex items-center gap-3 p-3 rounded-apple bg-nebula-50 dark:bg-nebula-800 cursor-pointer hover:bg-nebula-100 dark:hover:bg-nebula-700 transition-all">
                <item.icon className="w-5 h-5 text-[var(--text-secondary)]" />
                <div>
                  <p className="text-sm font-medium text-[var(--text-primary)]">{item.platform}</p>
                  <p className="text-xs text-[var(--text-secondary)]">{item.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </AppShell>
  )
}
