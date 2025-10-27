// Event Platform types
export interface Event {
  id: string
  organizerId: string
  title: string
  description: string
  codeOfConduct: string
  location: EventLocation
  dateTime: Date
  capacity: number
  ticketPrice: number
  status: EventStatus
  attendees: EventAttendee[]
  createdAt: Date
}

export interface EventLocation {
  type: 'physical' | 'virtual'
  address?: string
  city?: string
  state?: string
  country?: string
  virtualUrl?: string
}

export interface EventAttendee {
  userId: string
  registeredAt: Date
  consentSignature: DigitalConsent
  paymentStatus: PaymentStatus
}

export interface DigitalConsent {
  worldIdSignature: string
  codeOfConductHash: string
  signedAt: Date
  ipAddress: string
}

export type EventStatus = 'draft' | 'published' | 'cancelled' | 'completed'
export type PaymentStatus = 'pending' | 'completed' | 'failed' | 'refunded'