# Cursor Plan — Mighty Eagle

Deliver an **integratable paid MVP** for adult platforms: persona verification + consent proof + basic reputation + audit exports.

## North-star metrics
- Trial → paid conversion: ≥ 3%
- Weekly active API tenants: ≥ 10 (early)
- Verification success rate: ≥ 95% (excluding user failures)
- Webhook delivery success: ≥ 99% with retries

## Milestones
### M0 — Foundations (Week 1)
- Multi-tenant schema + API key auth
- Baseline event log (append-only)
- Minimal OpenAPI spec + SDK stub

### M1 — Persona verification (Week 2)
- Provider interface abstraction: `Verify()` → result + confidence
- Endpoints:
  - `POST /v1/persona/verifications`
  - `GET /v1/persona/verifications/{id}`
- Webhook events: `persona.verified`, `persona.failed`

### M2 — Consent tokens (Week 3)
- Consent issuance: `POST /v1/consent/tokens`
- Consent revocation: `POST /v1/consent/tokens/{id}/revoke`
- Receipt generation (MVP: server-signed payload; upgrade later to hash-chain)

### M3 — Reputation v1 (Week 4)
- Deterministic scoring:
  - verified status weight
  - account age
  - dispute signals (optional inputs)
- Endpoint: `GET /v1/reputation/{subject}`

### M4 — Audit exports + Webhooks reliability (Week 5)
- Export job:
  - `POST /v1/audit/exports`
  - `GET /v1/audit/exports/{id}`
- Webhook retries with backoff, DLQ table

### M5 — Billing & Entitlements (Week 6)
- Subscription + usage tier metering
- Tenant entitlements (caps):
  - monthly verifications
  - export frequency
- Optional Admin console: tenant usage + webhook health

### M6 — Hardening + Launch (Weeks 7–8)
- Threat model doc + security review
- Rate limiting + abuse prevention
- Privacy: retention policies + deletion requests
- Integration guide + sample app

## API (MVP)
- Auth: `X-API-Key` per tenant
- Persona:
  - `POST /v1/persona/verifications`
  - `GET /v1/persona/verifications/{id}`
- Consent:
  - `POST /v1/consent/tokens`
  - `POST /v1/consent/tokens/{id}/revoke`
- Reputation:
  - `GET /v1/reputation/{subject}`
- Webhooks:
  - `POST /v1/webhooks/endpoints`
  - `POST /v1/webhooks/endpoints/{id}/rotate-secret`
- Audit:
  - `POST /v1/audit/exports`
  - `GET /v1/audit/exports/{id}`

## Definition of done
- A platform can integrate in < 2 hours using API docs + sample app.
- Consent issuance/revocation is auditable and exportable.
- Billing gates work; usage limits enforced.
