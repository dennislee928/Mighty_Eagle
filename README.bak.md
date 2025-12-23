
# Aegis Trust Ecosystem

A Web3-driven trust ecosystem for adult communities, built with Turborepo, Next.js, and Go.

## Project Structure

```
aegis-trust-ecosystem/
├── apps/
│   └── web/                 # Next.js frontend application
├── services/
│   └── api-go/             # Go backend API service
├── packages/
│   ├── shared/             # Shared TypeScript types and utilities
│   ├── typescript-config/  # Shared TypeScript configurations
│   └── eslint-config/      # Shared ESLint configurations
└── scripts/                # Database and deployment scripts
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Go 1.21+
- Docker and Docker Compose

### Setup

1. **Clone and install dependencies:**
   ```bash
   git clone <repository-url>
   cd aegis-trust-ecosystem
   npm install
   ```

2. **Start the database services:**
   ```bash
   docker-compose up -d
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   cp services/api-go/.env.example services/api-go/.env
   ```

4. **Start the development servers:**
   ```bash
   npm run dev
   ```

This will start:
- Next.js web app on http://localhost:3000
- Go API service on http://localhost:8080
- PostgreSQL on localhost:5432
- Redis on localhost:6379

### Available Scripts

- `npm run dev` - Start all services in development mode
- `npm run build` - Build all applications
- `npm run lint` - Lint all code
- `npm run format` - Format code with Prettier
- `npm run type-check` - Run TypeScript type checking

### Database

The PostgreSQL database is automatically initialized with basic tables when you run `docker-compose up`. The database includes:

- Users table with World ID verification support
- UUID extension for primary keys
- Basic indexes for performance

### Services

#### Web App (`apps/web`)
- Next.js 14 with App Router
- TypeScript configuration
- ESLint and Prettier setup

#### API Service (`services/api-go`)
- Go with Gin framework
- PostgreSQL integration with GORM
- Redis integration for caching
- Health check endpoint at `/health`

#### Shared Packages
- `@aegis/shared` - Common TypeScript types and interfaces
- `@aegis/typescript-config` - Shared TypeScript configurations
- `@aegis/eslint-config` - Shared ESLint rules
