import api from '@/lib/api'
import type { Category, PaginatedProducts, Product, ListProductsParams } from '@/types'

export const productService = {
  list: async (params: ListProductsParams = {}): Promise<PaginatedProducts> => {
    const { data } = await api.get<{ success: boolean; data: PaginatedProducts }>(
      '/products',
      { params }
    )
    return data.data!
  },

  getById: async (id: string): Promise<Product> => {
    const { data } = await api.get<{ success: boolean; data: Product }>(`/products/${id}`)
    return data.data!
  },

  listCategories: async (): Promise<Category[]> => {
    const { data } = await api.get<{ success: boolean; data: Category[] }>('/categories')
    return data.data ?? []
  },
}
