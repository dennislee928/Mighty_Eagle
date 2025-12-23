-- Mighty Eagle Trust Layer - Initial Schema Migration
-- Milestone M0: Foundations

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- TENANTS & AUTHENTICATION
-- ============================================================================

-- Tenants table - Multi-tenant workspace management
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    tier VARCHAR(50) NOT NULL DEFAULT 'lite' CHECK (tier IN ('lite', 'pro', 'enterprise')),
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'deleted')),
    api_key VARCHAR(255) NOT NULL UNIQUE,
    api_secret_hash VARCHAR(255) NOT NULL, -- For webhook signing
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tenants_api_key ON tenants(api_key);
CREATE INDEX idx_tenants_status ON tenants(status);

-- ============================================================================
-- EVENT LOG (Append-only audit trail)
-- ============================================================================

CREATE TABLE event_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL, -- e.g., 'persona.verified', 'consent.issued'
    event_version VARCHAR(10) NOT NULL DEFAULT 'v1',
    actor_id VARCHAR(255), -- User/system that triggered the event
    subject_id VARCHAR(255), -- Entity the event is about
    resource_type VARCHAR(100), -- e.g., 'verification', 'consent_token'
    resource_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_event_log_tenant ON event_log(tenant_id, created_at DESC);
CREATE INDEX idx_event_log_type ON event_log(event_type, created_at DESC);
CREATE INDEX idx_event_log_subject ON event_log(subject_id);
CREATE INDEX idx_event_log_resource ON event_log(resource_type, resource_id);

-- ============================================================================
-- PERSONA VERIFICATION (M1)
-- ============================================================================

CREATE TABLE persona_verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    subject_id VARCHAR(255) NOT NULL, -- User identifier from tenant's system
    provider VARCHAR(50) NOT NULL, -- 'worldid', 'mock', etc.
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'failed', 'expired')),
    confidence_score DECIMAL(5,2), -- 0.00 to 100.00
    verification_data JSONB NOT NULL DEFAULT '{}', -- Provider-specific data
    proof_hash VARCHAR(255), -- Hash of verification proof
    expires_at TIMESTAMP WITH TIME ZONE,
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_persona_tenant_subject ON persona_verifications(tenant_id, subject_id);
CREATE INDEX idx_persona_status ON persona_verifications(status);
CREATE INDEX idx_persona_created ON persona_verifications(created_at DESC);

-- ============================================================================
-- CONSENT TOKENS (M2)
-- ============================================================================

CREATE TABLE consent_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    parties TEXT[] NOT NULL, -- Array of party identifiers
    scope VARCHAR(255) NOT NULL, -- e.g., 'messaging', 'content_access'
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    receipt_signature TEXT NOT NULL, -- Server-signed receipt
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'revoked', 'expired')),
    metadata JSONB DEFAULT '{}',
    issued_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    revoked_by VARCHAR(255),
    revoke_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_consent_tenant ON consent_tokens(tenant_id, created_at DESC);
CREATE INDEX idx_consent_status ON consent_tokens(status);
CREATE INDEX idx_consent_parties ON consent_tokens USING GIN(parties);
CREATE INDEX idx_consent_expires ON consent_tokens(expires_at);

-- ============================================================================
-- REPUTATION SCORES (M3)
-- ============================================================================

CREATE TABLE reputation_scores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    subject_id VARCHAR(255) NOT NULL,
    score DECIMAL(5,2) NOT NULL CHECK (score >= 0 AND score <= 100),
    verified_status_weight DECIMAL(5,2) DEFAULT 0,
    account_age_weight DECIMAL(5,2) DEFAULT 0,
    dispute_signals_weight DECIMAL(5,2) DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE, -- For caching
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, subject_id)
);

CREATE INDEX idx_reputation_tenant_subject ON reputation_scores(tenant_id, subject_id);
CREATE INDEX idx_reputation_score ON reputation_scores(score DESC);

-- ============================================================================
-- WEBHOOKS (M1, M4)
-- ============================================================================

CREATE TABLE webhook_endpoints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    secret VARCHAR(255) NOT NULL, -- For HMAC signature
    enabled BOOLEAN DEFAULT true,
    events TEXT[] NOT NULL, -- Array of event types to subscribe to
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_webhook_endpoints_tenant ON webhook_endpoints(tenant_id);
CREATE INDEX idx_webhook_endpoints_enabled ON webhook_endpoints(enabled);

CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    webhook_endpoint_id UUID NOT NULL REFERENCES webhook_endpoints(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES event_log(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'dlq')),
    attempt_count INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    request_payload JSONB NOT NULL,
    request_headers JSONB,
    response_status INTEGER,
    response_body TEXT,
    error_message TEXT,
    next_retry_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_webhook_deliveries_endpoint ON webhook_deliveries(webhook_endpoint_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_retry ON webhook_deliveries(next_retry_at) WHERE status = 'pending';

-- ============================================================================
-- AUDIT EXPORTS (M4)
-- ============================================================================

CREATE TABLE audit_export_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    format VARCHAR(20) NOT NULL CHECK (format IN ('csv', 'json')),
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    event_types TEXT[], -- Filter by event types
    download_url TEXT,
    file_size_bytes BIGINT,
    record_count INTEGER,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE -- Download link expiry
);

CREATE INDEX idx_audit_exports_tenant ON audit_export_jobs(tenant_id, created_at DESC);
CREATE INDEX idx_audit_exports_status ON audit_export_jobs(status);

-- ============================================================================
-- BILLING & ENTITLEMENTS (M5)
-- ============================================================================

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE UNIQUE,
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    tier VARCHAR(50) NOT NULL DEFAULT 'lite' CHECK (tier IN ('lite', 'pro', 'enterprise')),
    status VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('trialing', 'active', 'past_due', 'canceled', 'unpaid')),
    current_period_start TIMESTAMP WITH TIME ZONE,
    current_period_end TIMESTAMP WITH TIME ZONE,
    cancel_at TIMESTAMP WITH TIME ZONE,
    canceled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_subscriptions_tenant ON subscriptions(tenant_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);

CREATE TABLE usage_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    verification_count INTEGER DEFAULT 0,
    consent_token_count INTEGER DEFAULT 0,
    export_count INTEGER DEFAULT 0,
    api_request_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, period_start)
);

CREATE INDEX idx_usage_metrics_tenant ON usage_metrics(tenant_id, period_start DESC);

-- ============================================================================
-- UPDATED_AT TRIGGERS
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers to relevant tables
CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_persona_verifications_updated_at BEFORE UPDATE ON persona_verifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_consent_tokens_updated_at BEFORE UPDATE ON consent_tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reputation_scores_updated_at BEFORE UPDATE ON reputation_scores
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhook_endpoints_updated_at BEFORE UPDATE ON webhook_endpoints
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhook_deliveries_updated_at BEFORE UPDATE ON webhook_deliveries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_usage_metrics_updated_at BEFORE UPDATE ON usage_metrics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED DATA (Development)
-- ============================================================================

-- Create a test tenant
INSERT INTO tenants (name, tier, api_key, api_secret_hash) VALUES
    ('Test Platform', 'pro', 'test_sk_' || encode(gen_random_bytes(32), 'hex'), encode(gen_random_bytes(32), 'hex'));
