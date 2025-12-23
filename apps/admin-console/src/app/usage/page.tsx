export default function UsagePage() {
  return (
    <div className="flex flex-col gap-8">
      <header>
        <h2 className="text-3xl font-bold mb-2">Billing & Usage</h2>
        <p className="text-slate-400">Monitor your quota and subscription limits.</p>
      </header>

      <div className="flex flex-col gap-8 max-w-4xl">
        <section className="p-8 rounded-2xl border border-slate-800 bg-slate-900/40">
          <div className="flex justify-between items-center mb-8">
            <h3 className="text-xl font-semibold">Current Month Usage</h3>
            <span className="text-sm px-3 py-1 rounded-full bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">December 2025</span>
          </div>

          <div className="space-y-8">
            <div>
              <div className="flex justify-between text-sm mb-2">
                <span className="text-slate-400">Persona Verifications</span>
                <span className="font-mono">1,284 / 10,000</span>
              </div>
              <div className="h-2 w-full bg-slate-800 rounded-full overflow-hidden">
                <div className="h-full bg-indigo-500" style={{ width: '12.8%' }}></div>
              </div>
            </div>

            <div>
              <div className="flex justify-between text-sm mb-2">
                <span className="text-slate-400">Audit Exports</span>
                <span className="font-mono">8 / 100</span>
              </div>
              <div className="h-2 w-full bg-slate-800 rounded-full overflow-hidden">
                <div className="h-full bg-cyan-500" style={{ width: '8%' }}></div>
              </div>
            </div>

            <div>
              <div className="flex justify-between text-sm mb-2">
                <span className="text-slate-400">API Requests</span>
                <span className="font-mono">45,201 / 1,000,000</span>
              </div>
              <div className="h-2 w-full bg-slate-800 rounded-full overflow-hidden">
                <div className="h-full bg-slate-400" style={{ width: '4.5%' }}></div>
              </div>
            </div>
          </div>
        </section>

        <section className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/40">
            <h3 className="font-semibold mb-4">Subscription Plan</h3>
            <div className="flex items-center gap-2 text-2xl font-bold mb-1">
              Pro Plan
              <span className="text-xs font-normal text-slate-500 line-through tracking-tighter">$49/mo</span>
              <span className="text-xs font-normal text-emerald-400">Active</span>
            </div>
            <p className="text-sm text-slate-500 mb-6">Your next billing date is Jan 12, 2026.</p>
            <button className="text-sm px-4 py-2 border border-slate-700 rounded-lg hover:bg-slate-800 transition-colors">Manage Stripe</button>
          </div>

          <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/40">
            <h3 className="font-semibold mb-4">API Keys</h3>
            <div className="bg-black/40 p-3 rounded border border-slate-800 font-mono text-sm text-slate-400 truncate mb-4">
              me_live_4a89...c1e2
            </div>
            <div className="flex gap-3">
              <button className="text-sm px-4 py-2 bg-slate-800 rounded-lg hover:bg-slate-700 transition-colors">Copy</button>
              <button className="text-sm px-4 py-2 text-rose-400 hover:bg-rose-500/10 rounded-lg transition-colors">Roll Key</button>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
