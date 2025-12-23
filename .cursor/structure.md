# Repository Structure — Mighty Eagle

Designed as an **API-first trust layer** with optional admin console.

## Top-level tree (proposed)

```
mighty-eagle-trust-consent/
  README.md
  structure.md
  cursor/
    plan.md

  apps/
    console/                  # Admin Console (optional in MVP)
      src/
        app/
        components/
        features/
          tenants/
          audit/
          policies/
          exports/
      package.json

  services/
    api/
      cmd/
      internal/
        auth/
        tenants/
        persona/              # PoP providers + verification workflow
        consent/              # token issuance / revocation / receipts
        reputation/           # scoring + rules
        webhooks/             # delivery + retry
        audit/                # event log + exports
        policies/             # decision engine (MVP: simple thresholds)
      migrations/
      openapi/
      Dockerfile

  infra/
    docker-compose.yml
    postgres/
    scripts/

  docs/
    PRD.md
    THREAT_MODEL.md
    PRIVACY.md
    COMPLIANCE.md
    API.md
    ROADMAP.md

  .env.example
  Makefile
  LICENSE
```

## Tenant model (MVP)
- **Tenant**: platform workspace
- **API key**: per tenant
- **Policy**: verification requirements + thresholds (MVP: basic)
- **Event log**: append-only audit events (tamper-evident later)

## Core entities
- **PersonaVerification**
- **ConsentToken** (+ revocation record)
- **ReputationScore**
- **WebhookEndpoint** (+ delivery attempts)
- **AuditExportJob**

## Non-goals (MVP)
- Content hosting or content moderation pipelines
- Payment processing
- Full graph reputation at scale (plan for v0.3+)
