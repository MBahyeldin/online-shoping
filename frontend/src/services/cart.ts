import api from '@/lib/api'
import type { Cart } from '@/types'

export const cartService = {
  getCart: async (): Promise<Cart> => {
    const { data } = await api.get<{ success: boolean; data: Cart }>('/cart')
    return data.data!
  },

  addItem: async (productId: string, quantity: number): Promise<Cart> => {
    const { data } = await api.post<{ success: boolean; data: Cart }>('/cart/items', {
      product_id: productId,
      quantity,
    })
    return data.data!
  },

  updateItem: async (itemId: string, quantity: number): Promise<Cart> => {
    const { data } = await api.put<{ success: boolean; data: Cart }>(
      `/cart/items/${itemId}`,
      { quantity }
    )
    return data.data!
  },

  removeItem: async (itemId: string): Promise<Cart> => {
    const { data } = await api.delete<{ success: boolean; data: Cart }>(
      `/cart/items/${itemId}`
    )
    return data.data!
  },

  clearCart: async (): Promise<void> => {
    await api.delete('/cart')
  },
}
