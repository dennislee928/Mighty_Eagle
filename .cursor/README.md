# Mighty Eagle
> Trust, Consent, and Persona infrastructure for **adult-content platforms**.

Mighty Eagle provides a privacy-minded trust layer to reduce fraud and abuse while supporting compliance workflows in adult-content ecosystems.

## Scope statement (important)
- This project does **not** distribute adult content.
- This project provides **infrastructure** used by adult platforms: **persona verification**, **consent proof**, **reputation**, and **audit**.
- Payment processing remains the responsibility of the integrating platform (Mighty Eagle is not a merchant-of-record).

## Target customers (ICP)
Adult platforms and creator ecosystems requiring:
- anti‑Sybil defenses (one human → one account)
- consent verification and revocation
- reputation scoring for trust decisions
- compliance/audit exports

## MVP capabilities
### 1) Persona verification
- Proof‑of‑Personhood (PoP) integration (pluggable providers)
- One-person-one-account signals
- Risk score output for policy gates (e.g., restrict features until verified)

### 2) Consent tokens
- Issue and revoke consent tokens tied to:
  - parties (creator / user)
  - scope (interaction type / permissions)
  - expiry / revocation
- Cryptographically verifiable receipts (tamper-evident)

### 3) Reputation
- Lightweight reputation model:
  - verified status weighting
  - anti-fraud heuristics
  - optional graph-based reputation later

### 4) Webhooks + audit exports
- Webhook notifications for verification/consent events
- Exportable audit logs (CSV/JSON) for compliance review

## Product shape
- **API-first SaaS** + optional **Admin Console**
- Integrations via:
  - REST API
  - JS SDK (web)
  - Webhooks

## Pricing (initial)
- **Lite**: base verification + consent tokens (monthly + usage tier)
- **Pro**: reputation + risk policies + higher limits
- **Enterprise**: RBAC, SSO, SLA, custom retention, audit packs

## Safety & compliance considerations
- Age-gating / age verification is **mandatory** for customers (integration supported).
- Data minimization: store only what is required for verification/audit.
- Retention controls: configurable retention window per tenant.
- Abuse prevention: rate limits, anomaly detection, and admin lockouts.

## Architecture (proposed)
- **API**: Go (Gin) or Node (Fastify)
- **DB**: Postgres (tenants, events, tokens)
- **Queue**: background jobs for exports, webhook retries
- **Secrets**: KMS / Vault
- **Deployment**: Docker Compose (dev) → Kubernetes/Cloud Run (prod)

## Local development (placeholder)
1. `cp .env.example .env`
2. `docker compose up -d`
3. `make dev`

## Repository layout
See: [`structure.md`](./structure.md)

## License
TBD (recommend dual license: open core + commercial terms for hosted enterprise modules).
