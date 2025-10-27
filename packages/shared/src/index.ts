// Core types for the Aegis Trust Ecosystem

export interface WorldIDVerification {
  worldId: string;
  nullifierHash: string;
  merkleRoot: string;
  proof: string;
  verificationLevel: 'orb' | 'device';
  verifiedAt: Date;
}

export interface User {
  id: string;
  worldIdVerification: WorldIDVerification;
  profile: UserProfile;
  reputation: ReputationScore;
  createdAt: Date;
  updatedAt: Date;
}

export interface UserProfile {
  displayName?: string;
  bio?: string;
  avatar?: string;
}

export interface ReputationScore {
  userId: string;
  trustBadges: number;
  verificationLevel: VerificationLevel;
  lastUpdated: Date;
}

export type VerificationLevel = 'orb' | 'device' | 'unverified';

export interface APIError {
  code: string;
  message: string;
  details?: Record<string, any>;
  timestamp: string;
  requestId: string;
}