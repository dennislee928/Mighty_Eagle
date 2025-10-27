-- Initialize Aegis Trust Ecosystem Database
-- This script sets up the basic database structure

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create basic tables structure (will be expanded in later tasks)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    world_id VARCHAR(255) UNIQUE NOT NULL,
    nullifier_hash VARCHAR(255) UNIQUE NOT NULL,
    verification_level VARCHAR(20) NOT NULL DEFAULT 'device',
    profile_data JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_world_id ON users(world_id);
CREATE INDEX IF NOT EXISTS idx_users_nullifier_hash ON users(nullifier_hash);

-- Insert a test user for development
INSERT INTO users (world_id, nullifier_hash, verification_level, profile_data) 
VALUES (
    'test_world_id_123',
    'test_nullifier_hash_456', 
    'device',
    '{"displayName": "Test User", "bio": "Development test user"}'
) ON CONFLICT (world_id) DO NOTHING;