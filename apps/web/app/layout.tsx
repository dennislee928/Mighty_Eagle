import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Aegis Trust Ecosystem',
  description: 'Web3-driven sex-positive trust ecosystem',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-background font-sans antialiased">
        {children}
      </body>
    </html>
  )
}