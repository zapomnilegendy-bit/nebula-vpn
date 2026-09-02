'use client'

import { useState } from 'react'
import { LayoutDashboard, Users, Server, CreditCard, Ticket, BarChart3, Shield, Settings, LogOut } from 'lucide-react'
import clsx from 'clsx'

const nav = [
  { icon: LayoutDashboard, label: 'Dashboard', href: '/dashboard' },
  { icon: Users, label: 'Users', href: '/users' },
  { icon: Server, label: 'Servers', href: '/servers' },
  { icon: CreditCard, label: 'Plans', href: '/plans' },
  { icon: Ticket, label: 'Promocodes', href: '/promocodes' },
  { icon: BarChart3, label: 'Analytics', href: '/analytics' },
]

export default function AdminShell({ children }: { children: React.ReactNode }) {
  const [active, setActive] = useState('Dashboard')

  return (
    <div className="flex h-screen">
      <aside className="w-60 bg-[var(--bg-primary)] border-r border-black/5 dark:border-white/5 flex flex-col">
        <div className="px-5 py-5 flex items-center gap-3 border-b border-black/5 dark:border-white/5">
          <div className="w-7 h-7 rounded-apple bg-accent flex items-center justify-center">
            <Shield className="w-3.5 h-3.5 text-white" />
          </div>
          <div>
            <span className="font-semibold text-sm">NebulaVPN</span>
            <span className="block text-[10px] text-[var(--text-secondary)] font-medium tracking-wide">ADMIN</span>
          </div>
        </div>

        <nav className="flex-1 px-3 py-3 space-y-0.5">
          {nav.map(({ icon: Icon, label, href }) => (
            <a key={label} href={href} onClick={(e) => { e.preventDefault(); setActive(label) }}
              className={clsx(
                'flex items-center gap-3 px-3 py-2 rounded-apple text-sm font-medium transition-all',
                active === label ? 'bg-nebula-100 dark:bg-nebula-800 text-[var(--text-primary)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-nebula-50 dark:hover:bg-nebula-900'
              )}>
              <Icon className="w-4 h-4" /> {label}
            </a>
          ))}
        </nav>

        <div className="p-3 border-t border-black/5 dark:border-white/5 space-y-0.5">
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-apple text-sm text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-nebula-50 dark:hover:bg-nebula-900 transition-all">
            <Settings className="w-4 h-4" /> Settings
          </a>
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-apple text-sm text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-all">
            <LogOut className="w-4 h-4" /> Sign Out
          </a>
        </div>
      </aside>

      <main className="flex-1 overflow-auto">
        <div className="max-w-7xl mx-auto px-8 py-8">{children}</div>
      </main>
    </div>
  )
}
