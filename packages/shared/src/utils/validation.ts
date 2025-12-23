// Basic validation utilities
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

export const isValidRating = (rating: number): boolean => {
  return rating >= 1 && rating <= 5 && Number.isInteger(rating)
}

export const isValidWorldId = (worldId: string): boolean => {
  // Basic validation - actual validation will be done by World ID API
  return worldId.length > 0 && worldId.length < 256
}

export const sanitizeString = (input: string): string => {
  return input.trim().replace(/[<>]/g, '')
}