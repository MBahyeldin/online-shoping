import api from '@/lib/api'
import type { CreateOrderPayload, Order, PaginatedOrders } from '@/types'

export const orderService = {
  create: async (payload: CreateOrderPayload): Promise<Order> => {
    const { data } = await api.post<{ success: boolean; data: Order }>('/orders', payload)
    return data.data!
  },

  list: async (page = 1, limit = 10): Promise<PaginatedOrders> => {
    const { data } = await api.get<{ success: boolean; data: PaginatedOrders }>('/orders', {
      params: { page, limit },
    })
    return data.data!
  },

  getById: async (id: string): Promise<Order> => {
    const { data } = await api.get<{ success: boolean; data: Order }>(`/orders/${id}`)
    return data.data!
  },
}
