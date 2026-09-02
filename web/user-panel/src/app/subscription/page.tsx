'use client'

import AppShell from '@/components/AppShell'
import { CreditCard, Check, ArrowRight } from 'lucide-react'

const plans = [
  {
    name: 'Free',
    price: '$0',
    period: 'forever',
    features: ['1 device', '50 Mbps max', 'WireGuard only', '3 locations', 'Community support'],
    cta: 'Current Plan',
    active: false,
  },
  {
    name: 'Pro',
    price: '$9.99',
    period: '/month',
    features: ['5 devices', '500 Mbps max', 'All protocols', '30+ locations', 'Priority support', 'Split tunneling', 'Kill switch'],
    cta: 'Upgrade to Pro',
    active: true,
    popular: true,
  },
  {
    name: 'Ultimate',
    price: '$19.99',
    period: '/month',
    features: ['Unlimited devices', 'No speed limit', 'All protocols + TUIC', '50+ locations', '24/7 priority support', 'Split tunneling', 'Kill switch', 'Dedicated IP', 'CDN relay'],
    cta: 'Upgrade to Ultimate',
    active: false,
  },
]

export default function SubscriptionPage() {
  return (
    <AppShell>
      <div className="space-y-8">
        <div>
          <h1 className="text-2xl font-bold text-[var(--text-primary)]">Subscription</h1>
          <p className="text-sm text-[var(--text-secondary)] mt-1">Manage your plan and billing.</p>
        </div>

        {/* Current plan card */}
        <div className="card p-6 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-apple bg-accent/10 flex items-center justify-center">
              <CreditCard className="w-6 h-6 text-accent" />
            </div>
            <div>
              <p className="text-sm text-[var(--text-secondary)]">Current Plan</p>
              <h2 className="text-lg font-bold text-[var(--text-primary)]">Pro · $9.99/month</h2>
              <p className="text-xs text-[var(--text-secondary)] mt-0.5">Renews on Jun 15, 2026 · 12 days remaining</p>
            </div>
          </div>
          <button className="btn-secondary text-xs">Manage Billing</button>
        </div>

        {/* Plan comparison */}
        <div className="grid grid-cols-3 gap-4">
          {plans.map((plan) => (
            <div key={plan.name} className={`card p-6 relative ${plan.popular ? 'ring-2 ring-accent' : ''}`}>
              {plan.popular && (
                <span className="absolute -top-2.5 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full bg-accent text-white text-[10px] font-bold tracking-wide uppercase">
                  Popular
                </span>
              )}
              <div className="text-center">
                <h3 className="text-sm font-semibold text-[var(--text-secondary)]">{plan.name}</h3>
                <div className="mt-3">
                  <span className="text-3xl font-bold text-[var(--text-primary)]">{plan.price}</span>
                  <span className="text-sm text-[var(--text-secondary)]">{plan.period}</span>
                </div>
              </div>

              <ul className="mt-6 space-y-3">
                {plan.features.map((f, i) => (
                  <li key={i} className="flex items-start gap-2 text-sm">
                    <Check className="w-4 h-4 text-[#34C759] mt-0.5 flex-shrink-0" />
                    <span className="text-[var(--text-primary)]">{f}</span>
                  </li>
                ))}
              </ul>

              <button className={`w-full mt-6 py-3 rounded-apple text-sm font-semibold transition-all ${
                plan.active
                  ? 'bg-nebula-100 dark:bg-nebula-800 text-[var(--text-primary)] cursor-default'
                  : 'btn-primary'
              }`}>
                {plan.cta}
              </button>
            </div>
          ))}
        </div>

        {/* Payment methods */}
        <div className="card p-6">
          <h3 className="text-sm font-semibold text-[var(--text-primary)] mb-4">Payment Method</h3>
          <div className="flex items-center gap-4 p-4 rounded-apple bg-nebula-50 dark:bg-nebula-800">
            <div className="w-10 h-7 rounded bg-blue-600 flex items-center justify-center text-white text-[10px] font-bold">
              VISA
            </div>
            <div className="flex-1">
              <p className="text-sm font-medium text-[var(--text-primary)]">Visa ending in 4242</p>
              <p className="text-xs text-[var(--text-secondary)]">Expires 12/28</p>
            </div>
            <button className="text-sm text-accent font-medium">Change</button>
          </div>
        </div>
      </div>
    </AppShell>
  )
}
