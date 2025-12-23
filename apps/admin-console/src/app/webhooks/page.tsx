export default function WebhooksPage() {
  return (
    <div className="flex flex-col gap-8">
      <header className="flex justify-between items-end">
        <div>
          <h2 className="text-3xl font-bold mb-2">Webhooks</h2>
          <p className="text-slate-400">Manage your real-time event notifications.</p>
        </div>
        <button className="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 rounded-lg font-medium transition-colors">
          Add Endpoint
        </button>
      </header>

      <div className="grid grid-cols-1 gap-6">
        <section className="overflow-hidden rounded-xl border border-slate-800 bg-slate-900/40">
          <table className="w-full text-left">
            <thead className="bg-slate-900/60 border-b border-slate-800">
              <tr>
                <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">URL</th>
                <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Events</th>
                <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Status</th>
                <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Last Delivery</th>
                <th className="px-6 py-4 text-sm font-semibold uppercase tracking-wider text-slate-400">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800">
              <tr>
                <td className="px-6 py-4">
                  <div className="font-medium">https://api.your-platform.com/webhooks</div>
                  <div className="text-xs text-slate-500">whsec_e39f...9c22</div>
                </td>
                <td className="px-6 py-4">
                  <div className="flex gap-1 flex-wrap">
                    <span className="text-[10px] px-2 py-0.5 rounded bg-slate-800 border border-slate-700">persona.*</span>
                    <span className="text-[10px] px-2 py-0.5 rounded bg-slate-800 border border-slate-700">consent.*</span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span className="flex items-center gap-1.5 text-xs text-emerald-400">
                    <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.5)]"></span>
                    Enabled
                  </span>
                </td>
                <td className="px-6 py-4">
                  <div className="text-sm">2 hours ago</div>
                  <div className="text-xs text-emerald-500">HTTP 200 OK</div>
                </td>
                <td className="px-6 py-4">
                  <button className="text-xs text-slate-400 hover:text-slate-100 uppercase tracking-tighter font-bold">Manage</button>
                </td>
              </tr>
              <tr>
                <td className="px-6 py-4">
                  <div className="font-medium text-slate-500 italic">https://staging.your-platform.com/hooks</div>
                </td>
                <td className="px-6 py-4 text-slate-500">*</td>
                <td className="px-6 py-4">
                  <span className="flex items-center gap-1.5 text-xs text-slate-500">
                    <span className="w-1.5 h-1.5 rounded-full bg-slate-600"></span>
                    Disabled
                  </span>
                </td>
                <td className="px-6 py-4 text-slate-500">—</td>
                <td className="px-6 py-4">
                  <button className="text-xs text-slate-400 hover:text-slate-100 uppercase tracking-tighter font-bold">Manage</button>
                </td>
              </tr>
            </tbody>
          </table>
        </section>

        <section className="p-6 rounded-xl border border-slate-800 bg-slate-900/40">
          <h3 className="text-lg font-semibold mb-6">Recent Failures</h3>
          <div className="flex flex-col gap-4">
            <div className="p-4 rounded-lg bg-rose-500/5 border border-rose-500/10 flex justify-between items-center">
              <div className="flex gap-4">
                <div className="w-10 h-10 rounded bg-rose-500/10 flex items-center justify-center font-bold text-rose-400">404</div>
                <div>
                  <div className="text-sm font-medium">Endpoint unreachable</div>
                  <div className="text-xs text-slate-500">persona.verified • 12 attempts • 1 hour ago</div>
                </div>
              </div>
              <button className="px-3 py-1.5 text-xs bg-rose-500/20 text-rose-400 rounded-md hover:bg-rose-500/30 transition-colors">Move to DLQ</button>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
