import { create } from 'zustand'
import type { Cart } from '@/types'

interface CartState {
  cart: Cart | null
  isOpen: boolean
  setCart: (cart: Cart) => void
  clearCart: () => void
  toggleCart: () => void
  openCart: () => void
  closeCart: () => void
  itemCount: () => number
}

export const useCartStore = create<CartState>()((set, get) => ({
  cart: null,
  isOpen: false,

  setCart: (cart) => set({ cart }),
  clearCart: () => set({ cart: null }),
  toggleCart: () => set((s) => ({ isOpen: !s.isOpen })),
  openCart: () => set({ isOpen: true }),
  closeCart: () => set({ isOpen: false }),

  itemCount: () => {
    const { cart } = get()
    if (!cart) return 0
    return cart.items.reduce((sum, item) => sum + item.quantity, 0)
  },
}))
