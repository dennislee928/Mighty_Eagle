import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Mighty Eagle Admin",
  description: "Trust Layer Administration",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className} bg-slate-950 text-slate-50 min-h-screen`}>
        <div className="flex">
          <aside className="w-64 border-r border-slate-800 h-screen sticky top-0 p-6 flex flex-col gap-6">
            <h1 className="text-xl font-bold bg-gradient-to-r from-indigo-400 to-cyan-400 bg-clip-text text-transparent">
              Mighty Eagle
            </h1>
            <nav className="flex flex-col gap-1">
              <a href="/" className="px-4 py-2 rounded-md hover:bg-slate-900 transition-colors">Dashboard</a>
              <a href="/verifications" className="px-4 py-2 rounded-md hover:bg-slate-900 transition-colors">Verifications</a>
              <a href="/webhooks" className="px-4 py-2 rounded-md hover:bg-slate-900 transition-colors">Webhooks</a>
              <a href="/usage" className="px-4 py-2 rounded-md hover:bg-slate-900 transition-colors">Billing & Usage</a>
            </nav>
            <div className="mt-auto pt-6 border-t border-slate-800 text-sm text-slate-500">
              v0.1.0-alpha
            </div>
          </aside>
          <main className="flex-1 p-8">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
