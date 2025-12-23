export default function Dashboard() {
  return (
    <div className="flex flex-col gap-8">
      <header>
        <h2 className="text-3xl font-bold mb-2">Dashboard</h2>
        <p className="text-slate-400">Welcome to your Trust Layer control center.</p>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/50 backdrop-blur-sm">
          <h3 className="text-slate-400 text-sm font-medium mb-4 uppercase tracking-wider">Total Verifications</h3>
          <div className="text-3xl font-bold font-mono">1,284</div>
          <div className="text-emerald-400 text-xs mt-2">↑ 12% from last month</div>
        </div>
        <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/50 backdrop-blur-sm">
          <h3 className="text-slate-400 text-sm font-medium mb-4 uppercase tracking-wider">Health</h3>
          <div className="text-3xl font-bold text-emerald-500">Normal</div>
          <div className="text-slate-500 text-xs mt-2">All systems operational</div>
        </div>
        <div className="p-6 rounded-xl border border-slate-800 bg-slate-900/50 backdrop-blur-sm">
          <h3 className="text-slate-400 text-sm font-medium mb-4 uppercase tracking-wider">Tier</h3>
          <div className="text-3xl font-bold">Pro Plan</div>
          <div className="text-indigo-400 text-xs mt-2">Manage subscription →</div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <section className="p-6 rounded-xl border border-slate-800 bg-slate-900/50">
          <h3 className="text-lg font-semibold mb-6">Recent Verifications</h3>
          <div className="flex flex-col gap-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="flex items-center justify-between p-3 rounded-lg bg-slate-800/30">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-slate-700 flex items-center justify-center text-xs">US</div>
                  <div>
                    <div className="font-medium text-sm">user_{i * 123}</div>
                    <div className="text-xs text-slate-500">World ID • 2m ago</div>
                  </div>
                </div>
                <div className="text-xs px-2 py-1 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/20">Verified</div>
              </div>
            ))}
          </div>
        </section>

        <section className="p-6 rounded-xl border border-slate-800 bg-slate-900/50">
          <h3 className="text-lg font-semibold mb-6">Webhook Health</h3>
          <div className="h-48 flex items-end gap-2 overflow-hidden px-2">
            {[40, 60, 30, 80, 70, 90, 100, 95, 80, 85, 90, 100, 100, 98, 100].map((h, i) => (
              <div 
                key={i} 
                className="flex-1 bg-indigo-500/40 rounded-t-sm" 
                style={{ height: `${h}%` }}
              ></div>
            ))}
          </div>
          <div className="flex justify-between mt-4 text-xs text-slate-500">
            <span>24h ago</span>
            <span>Now</span>
          </div>
        </section>
      </div>
    </div>
  );
}
