import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Filter, SortAsc, SortDesc, ChevronLeft, ChevronRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ProductCard } from '@/components/shared/ProductCard'
import { PageLoader } from '@/components/shared/LoadingSpinner'
import { productService } from '@/services/products'
import type { ListProductsParams } from '@/types'

export function ProductsPage() {
  const [params, setParams] = useState<ListProductsParams>({
    page: 1,
    limit: 12,
  })

  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: productService.listCategories,
    staleTime: Infinity,
  })

  const { data, isLoading, isFetching } = useQuery({
    queryKey: ['products', params],
    queryFn: () => productService.list(params),
    placeholderData: (prev) => prev,
  })

  const setCategory = (value: string) =>
    setParams((p) => ({ ...p, category_id: value === 'all' ? undefined : value, page: 1 }))

  const setSort = (value: string) =>
    setParams((p) => ({
      ...p,
      sort: value === 'none' ? undefined : (value as ListProductsParams['sort']),
    }))

  const setPage = (page: number) => setParams((p) => ({ ...p, page }))

  const totalPages = data?.total_pages ?? 1

  return (
    <main className="min-h-screen">
      {/* Page header */}
      <div className="bg-gradient-to-br from-orange-50 to-amber-50 border-b py-12">
        <div className="container mx-auto max-w-6xl px-4">
          <h1 className="text-4xl font-extrabold text-gray-900">Our Cakes</h1>
          <p className="mt-2 text-muted-foreground">
            {data ? `${data.total} cakes available` : 'Browse our full collection'}
          </p>
        </div>
      </div>

      <div className="container mx-auto max-w-6xl px-4 py-8">
        {/* Filters bar */}
        <div className="flex flex-wrap gap-4 items-center mb-8">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Filter className="h-4 w-4" />
            <span className="font-medium">Filters:</span>
          </div>

          {/* Category filter */}
          <Select
            onValueChange={setCategory}
            defaultValue="all"
          >
            <SelectTrigger className="w-48">
              <SelectValue placeholder="All Categories" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Categories</SelectItem>
              {categories?.map((cat) => (
                <SelectItem key={cat.id} value={cat.id}>
                  {cat.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          {/* Sort */}
          <Select onValueChange={setSort} defaultValue="none">
            <SelectTrigger className="w-44">
              <SelectValue placeholder="Sort By" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">Default</SelectItem>
              <SelectItem value="price_asc">
                <span className="flex items-center gap-2">
                  <SortAsc className="h-4 w-4" /> Price: Low to High
                </span>
              </SelectItem>
              <SelectItem value="price_desc">
                <span className="flex items-center gap-2">
                  <SortDesc className="h-4 w-4" /> Price: High to Low
                </span>
              </SelectItem>
            </SelectContent>
          </Select>

          {isFetching && !isLoading && (
            <span className="text-xs text-muted-foreground animate-pulse">Updating…</span>
          )}
        </div>

        {/* Products grid */}
        {isLoading ? (
          <PageLoader />
        ) : data?.products.length === 0 ? (
          <div className="py-20 text-center">
            <p className="text-4xl mb-4">🎂</p>
            <p className="text-lg font-semibold text-gray-900">No cakes found</p>
            <p className="text-muted-foreground mt-2">
              Try adjusting your filters or check back later.
            </p>
            <Button
              variant="outline"
              className="mt-4"
              onClick={() => setParams({ page: 1, limit: 12 })}
            >
              Clear Filters
            </Button>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {data?.products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-2 mt-10">
                <Button
                  variant="outline"
                  size="icon"
                  disabled={params.page === 1}
                  onClick={() => setPage((params.page ?? 1) - 1)}
                >
                  <ChevronLeft className="h-4 w-4" />
                </Button>

                {Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
                  <Button
                    key={p}
                    variant={p === params.page ? 'default' : 'outline'}
                    size="icon"
                    onClick={() => setPage(p)}
                    className="h-9 w-9"
                  >
                    {p}
                  </Button>
                ))}

                <Button
                  variant="outline"
                  size="icon"
                  disabled={params.page === totalPages}
                  onClick={() => setPage((params.page ?? 1) + 1)}
                >
                  <ChevronRight className="h-4 w-4" />
                </Button>
              </div>
            )}
          </>
        )}
      </div>
    </main>
  )
}
