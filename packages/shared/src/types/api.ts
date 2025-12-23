// API response types
export interface APIResponse<T = any> {
  success: boolean
  data?: T
  error?: APIError
  timestamp: string
}

export interface APIError {
  code: string
  message: string
  details?: Record<string, any>
  timestamp: string
  requestId: string
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  limit: number
  hasNext: boolean
  hasPrev: boolean
}

// Common API request types
export interface CreateReviewRequest {
  productId: string
  rating: number
  title: string
  content: string
}

export interface CreateEventRequest {
  title: string
  description: string
  codeOfConduct: string
  location: EventLocation
  dateTime: string
  capacity: number
  ticketPrice: number
}

export interface CreateConsentualLinkRequest {
  recipientId: string
  eventId?: string
}