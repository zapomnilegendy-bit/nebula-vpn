'use client'

import { useState } from 'react'
import { Shield, Globe, Monitor, CreditCard, ChevronRight, Power, Settings, LogOut } from 'lucide-react'
import clsx from 'clsx'

const navItems = [
  { icon: Shield, label: 'Dashboard', href: '/dashboard' },
  { icon: Globe, label: 'Servers', href: '/servers' },
  { icon: Monitor, label: 'Devices', href: '/devices' },
  { icon: CreditCard, label: 'Subscription', href: '/subscription' },
]

export default function AppShell({ children }: { children: React.ReactNode }) {
  const [active, setActive] = useState('Dashboard')
  const [isConnected, setIsConnected] = useState(false)

  return (
    <div className="flex h-screen bg-[var(--bg-secondary)]">
      {/* Sidebar */}
      <aside className="w-64 bg-[var(--bg-primary)] border-r border-black/5 dark:border-white/5 flex flex-col">
        {/* Logo */}
        <div className="px-6 py-5 flex items-center gap-3">
          <div className="w-8 h-8 rounded-apple bg-black dark:bg-white flex items-center justify-center">
            <Shield className="w-4 h-4 text-white dark:text-black" />
          </div>
          <span className="font-semibold text-base">NebulaVPN</span>
        </div>

        {/* Nav */}
        <nav className="flex-1 px-3 py-2 space-y-1">
          {navItems.map(({ icon: Icon, label, href }) => (
            <a
              key={label}
              href={href}
              onClick={(e) => { e.preventDefault(); setActive(label) }}
              className={clsx(
                'flex items-center gap-3 px-3 py-2.5 rounded-apple text-sm font-medium transition-all duration-150',
                active === label
                  ? 'bg-nebula-100 dark:bg-nebula-800 text-[var(--text-primary)]'
                  : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-nebula-50 dark:hover:bg-nebula-900'
              )}
            >
              <Icon className="w-4 h-4" />
              {label}
            </a>
          ))}
        </nav>

        {/* Bottom */}
        <div className="p-3 border-t border-black/5 dark:border-white/5 space-y-1">
          <a href="#" className="flex items-center gap-3 px-3 py-2.5 rounded-apple text-sm font-medium text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-nebula-50 dark:hover:bg-nebula-900 transition-all">
            <Settings className="w-4 h-4" />
            Settings
          </a>
          <a href="#" className="flex items-center gap-3 px-3 py-2.5 rounded-apple text-sm font-medium text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-all">
            <LogOut className="w-4 h-4" />
            Sign Out
          </a>
        </div>
      </aside>

      {/* Main */}
      <main className="flex-1 overflow-auto">
        <div className="max-w-5xl mx-auto px-8 py-8">
          {children}
        </div>
      </main>
    </div>
  )
}
