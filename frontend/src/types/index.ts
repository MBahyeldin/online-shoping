// ─── API Envelope ─────────────────────────────────────────────────────────────
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
}

// ─── Auth ────────────────────────────────────────────────────────────────────
export interface User {
  id: string
  first_name: string
  last_name: string
  email: string
  phone: string
}

export interface AuthResult {
  token: string
  user: User
}

export interface RegisterPayload {
  first_name: string
  last_name: string
  phone_number: string
  email: string
}

export interface VerifyOTPPayload {
  email: string
  otp: string
}

// ─── Category ─────────────────────────────────────────────────────────────────
export interface Category {
  id: string
  name: string
  slug: string
  created_at: string
  updated_at: string
}

// ─── Product ──────────────────────────────────────────────────────────────────
export interface Product {
  id: string
  category_id: string | null
  category_name: string | null
  category_slug: string | null
  name: string
  description: string | null
  price: number
  image_url: string | null
  stock_quantity: number
  is_active: boolean
}

export interface ListProductsParams {
  page?: number
  limit?: number
  category_id?: string
  sort?: 'price_asc' | 'price_desc'
}

export interface PaginatedProducts {
  products: Product[]
  total: number
  page: number
  limit: number
  total_pages: number
}

// ─── Cart ────────────────────────────────────────────────────────────────────
export interface CartItem {
  id: string
  product_id: string
  product_name: string
  product_image_url: string | null
  price: number
  quantity: number
  subtotal: number
}

export interface Cart {
  id: string
  items: CartItem[]
  total: number
}

// ─── Order ───────────────────────────────────────────────────────────────────
export interface OrderItem {
  id: string
  product_id: string
  product_name: string
  image_url: string | null
  quantity: number
  unit_price: number
  total_price: number
}

export interface Order {
  id: string
  delivery_address: string
  delivery_date: string
  notes: string | null
  payment_method: string
  status: 'pending' | 'confirmed' | 'preparing' | 'delivered' | 'cancelled'
  total_amount: number
  items: OrderItem[]
  created_at: string
}

export interface CreateOrderPayload {
  delivery_address: string
  delivery_date: string // RFC3339
  notes?: string
  payment_method: string
}

export interface PaginatedOrders {
  orders: Order[]
  total: number
  page: number
  limit: number
  total_pages: number
}
