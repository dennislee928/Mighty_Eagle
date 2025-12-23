export default function VerificationsPage() {
  return (
    <div className="flex flex-col gap-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold mb-2">Verifications</h2>
          <p className="text-slate-400">View and manage all persona verification records.</p>
        </div>
        <div className="flex gap-3">
          <input 
            type="text" 
            placeholder="Search Subject ID..." 
            className="px-4 py-2 bg-slate-900 border border-slate-800 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-indigo-500 w-64"
          />
          <button className="px-4 py-2 border border-slate-700 rounded-lg text-sm hover:bg-slate-800 transition-colors">
            Filter
          </button>
        </div>
      </header>

      <div className="overflow-hidden rounded-xl border border-slate-800 bg-slate-900/40">
        <table className="w-full text-left">
          <thead className="bg-slate-900/60 border-b border-slate-800">
            <tr>
              <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">ID</th>
              <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Subject</th>
              <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Provider</th>
              <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Status</th>
              <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Confidence</th>
              <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Time</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800">
            {[1, 2, 3, 4, 5, 6, 7, 8].map((i) => (
              <tr key={i} className="hover:bg-slate-900/20 transition-colors">
                <td className="px-6 py-4 font-mono text-xs text-slate-500">
                  {Math.random().toString(16).substring(2, 10)}...
                </td>
                <td className="px-6 py-4 font-medium">user_ox_{i * 999}</td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2">
                    <span className="w-2 h-2 rounded-full bg-indigo-500"></span>
                    World ID
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span className="px-2 py-0.5 rounded text-[10px] bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 uppercase font-bold">
                    Verified
                  </span>
                </td>
                <td className="px-6 py-4 font-mono">0.98</td>
                <td className="px-6 py-4 text-sm text-slate-500">12/23/2025 10:45 AM</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
