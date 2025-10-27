// Reputation System types
export interface ConsentualLink {
  id: string
  initiatorId: string
  recipientId: string
  eventId?: string
  status: LinkStatus
  initiatorSignature?: WorldIDSignature
  recipientSignature?: WorldIDSignature
  createdAt: Date
  completedAt?: Date
}

export interface TrustBadge {
  id: string
  userId: string
  linkedUserId: string
  eventContext?: string
  awardedAt: Date
}

export interface WorldIDSignature {
  signature: string
  timestamp: Date
  nullifierHash: string
}

export type LinkStatus = 'pending' | 'accepted' | 'completed' | 'expired' | 'declined'