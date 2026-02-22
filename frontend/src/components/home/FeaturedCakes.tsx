import { Link } from 'react-router-dom'
import { ArrowRight } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { productService } from '@/services/products'
import { ProductCard } from '@/components/shared/ProductCard'
import { PageLoader } from '@/components/shared/LoadingSpinner'
import { Button } from '@/components/ui/button'

export function FeaturedCakes() {
  const { data, isLoading } = useQuery({
    queryKey: ['products', 'featured'],
    queryFn: () => productService.list({ page: 1, limit: 4 }),
  })

  return (
    <section className="py-20 bg-white">
      <div className="container mx-auto max-w-6xl px-4">
        {/* Section header */}
        <div className="text-center mb-12">
          <p className="text-sm font-semibold text-primary uppercase tracking-widest mb-2">
            Our Specialties
          </p>
          <h2 className="text-4xl font-extrabold text-gray-900">Featured Cakes</h2>
          <p className="mt-3 text-muted-foreground max-w-xl mx-auto">
            Handcrafted by our expert bakers using only the finest ingredients. Browse a
            selection of our most popular creations.
          </p>
        </div>

        {/* Products grid */}
        {isLoading ? (
          <PageLoader />
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {(data?.products ?? []).map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        )}

        {/* CTA */}
        <div className="mt-10 text-center">
          <Button variant="outline" size="lg" asChild>
            <Link to="/products">
              View All Cakes
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </div>
      </div>
    </section>
  )
}
