# Mighty Eagle - Developer Documentation

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (for local development)

### Local Development Setup

1. **Clone the repository**
```bash
git clone https://github.com/dennislee928/mighty-eagle.git
cd mighty-eagle
```

2. **Set up environment**
```bash
cp services/api-go/.env.example services/api-go/.env
# Edit .env with your local configuration
```

3. **Start infrastructure with Docker Compose**
```bash
docker compose up -d postgres redis
```

4. **Run database migrations**
```bash
make migrate-up
```

5. **Install dependencies**
```bash
cd services/api-go
go mod download
```

6. **Start the API server**
```bash
make dev
# Or: cd services/api-go && go run main.go
```

The API will be available at `http://localhost:8080`

### Creating a Test Tenant

```bash
# Connect to database
psql -h localhost -U postgres -d mighty_eagle

# Create a test tenant
INSERT INTO tenants (name, tier, api_key, api_secret_hash) 
VALUES (
  'My Test Platform', 
  'pro', 
  'me_pro_sk_' || encode(gen_random_bytes(32), 'hex'),
  encode(gen_random_bytes(32), 'hex')
) 
RETURNING id, api_key;
```

### Testing the API

```bash
# Health check (no auth required)
curl http://localhost:8080/health

# API info (with auth)
curl http://localhost:8080/v1/info \
  -H "X-API-Key: YOUR_API_KEY_HERE"
```

## Project Structure

```
mighty-eagle/
├── services/
│   └── api-go/                   # Go API service
│       ├── cmd/                  # Command-line tools
│       ├── config/               # Configuration (DB, Redis)
│       ├── internal/
│       │   ├── audit/            # Event logging
│       │   ├── middleware/       # Auth, rate limiting, CORS
│       │   ├── models/           # GORM models
│       │   ├── persona/          # Persona verification (M1)
│       │   ├── consent/          # Consent tokens (M2)
│       │   ├── reputation/       # Reputation scoring (M3)
│       │   ├── webhooks/         # Webhook delivery (M4)
│       │   ├── tenants/          # Tenant management
│       │   └── router/           # Route setup
│       ├── migrations/           # Database migrations
│       ├── openapi/              # OpenAPI spec
│       ├── main.go
│       ├── Dockerfile
│       └── go.mod
├── apps/
│   └── console/                  # Admin console (M5)
├── packages/                     # Shared packages
├── docker-compose.yml
├── Makefile
└── README.md
```

## API Documentation

### Authentication

All API v1 endpoints require an API key in the `X-API-Key` header:

```bash
X-API-Key: me_pro_sk_abc123...
```

### Rate Limits

Rate limits are tier-based:
- **Lite**: 60 requests/minute
- **Pro**: 600 requests/minute
- **Enterprise**: 6000 requests/minute

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Total requests allowed per minute
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Unix timestamp when limit resets

### Available Endpoints

For full API documentation, see the [OpenAPI spec](./services/api-go/openapi/openapi.yaml)

**M0 - Foundations (Completed)**
- `GET /health` - Health check
- `GET /v1/info` - API information

**M1 - Persona Verification (TODO)**
- `POST /v1/persona/verifications` - Create verification
- `GET /v1/persona/verifications/{id}` - Get verification status

**M2 - Consent Tokens (TODO)**
- `POST /v1/consent/tokens` - Issue consent token
- `POST /v1/consent/tokens/{id}/revoke` - Revoke token

**M3 - Reputation (TODO)**
- `GET /v1/reputation/{subject}` - Get reputation score

**M4 - Webhooks & Audit (TODO)**
- `POST /v1/webhooks/endpoints` - Register webhook
- `POST /v1/audit/exports` - Create audit export

## Development Commands

```bash
# Start development server
make dev

# Build production binary
make build

# Run tests
make test

# Format code
make fmt

# Database migrations
make migrate-up

# Docker commands
make docker-up        # Start all services
make docker-down      # Stop all services
make docker-logs      # View logs
make docker-rebuild   # Rebuild and restart
```

## Database Schema

See [migrations/001_initial_schema.sql](./services/api-go/migrations/001_initial_schema.sql) for complete schema.

### Core Tables
- `tenants` - Multi-tenant workspaces
- `event_log` - Append-only audit trail
- `persona_verifications` - Verification records
- `consent_tokens` - Consent management
- `reputation_scores` - Reputation data
- `webhook_endpoints` - Webhook configurations
- `webhook_deliveries` - Delivery tracking
- `audit_export_jobs` - Export jobs

## Next Steps

The M0 foundation is complete. Next milestones:

- **M1**: Implement persona verification with World ID integration
- **M2**: Build consent token issuance and revocation
- **M3**: Create reputation scoring algorithm
- **M4**: Add webhook delivery and audit exports
- **M5**: Build admin console and billing
- **M6**: Security hardening and launch prep

## Contributing

(TBD)

## License

(TBD)
