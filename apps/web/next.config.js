/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    transpilePackages: ["@aegis/ui", "@aegis/shared"]
  }
}

module.exports = nextConfig