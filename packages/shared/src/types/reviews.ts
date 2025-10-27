// Review Library types
export interface Product {
  id: string
  name: string
  category: ProductCategory
  description: string
  affiliateLinks: AffiliateLink[]
  averageRating: number
  reviewCount: number
  createdAt: Date
}

export interface Review {
  id: string
  productId: string
  userId: string
  rating: number
  title: string
  content: string
  isVerifiedHuman: boolean
  helpfulVotes: number
  createdAt: Date
}

export interface AffiliateLink {
  provider: string
  url: string
  commission: number
}

export type ProductCategory = 
  | 'toys'
  | 'wellness'
  | 'books'
  | 'accessories'
  | 'other'