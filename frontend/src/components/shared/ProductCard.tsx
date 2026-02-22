import { ShoppingCart, ImageOff } from 'lucide-react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/store/authStore'
import { useCartStore } from '@/store/cartStore'
import { cartService } from '@/services/cart'
import { formatCurrency } from '@/lib/utils'
import type { Product } from '@/types'

interface ProductCardProps {
  product: Product
}

export function ProductCard({ product }: ProductCardProps) {
  const { isAuthenticated } = useAuthStore()
  const { setCart, openCart } = useCartStore()
  const navigate = useNavigate()
  const [adding, setAdding] = useState(false)
  const [imgError, setImgError] = useState(false)

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      navigate('/register')
      return
    }

    setAdding(true)
    try {
      const updatedCart = await cartService.addItem(product.id, 1)
      setCart(updatedCart)
      openCart()
    } catch (err) {
      console.error('Add to cart failed:', err)
    } finally {
      setAdding(false)
    }
  }

  const outOfStock = product.stock_quantity === 0

  return (
    <div className="group flex flex-col rounded-xl border bg-card shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      {/* Image */}
      <div className="relative overflow-hidden aspect-[4/3] bg-muted">
        {product.image_url && !imgError ? (
          <img
            src={product.image_url}
            alt={product.name}
            className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-300"
            onError={() => setImgError(true)}
            loading="lazy"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center">
            <ImageOff className="h-12 w-12 text-muted-foreground/30" />
          </div>
        )}

        {/* Category badge */}
        {product.category_name && (
          <Badge className="absolute top-3 left-3 text-xs" variant="secondary">
            {product.category_name}
          </Badge>
        )}

        {/* Out of stock overlay */}
        {outOfStock && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/40">
            <Badge variant="destructive" className="text-sm px-3 py-1">
              Out of Stock
            </Badge>
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex flex-col flex-1 p-4 gap-3">
        <div className="flex-1">
          <h3 className="font-semibold text-base leading-tight line-clamp-2">{product.name}</h3>
          {product.description && (
            <p className="mt-1.5 text-sm text-muted-foreground line-clamp-2">
              {product.description}
            </p>
          )}
        </div>

        <div className="flex items-center justify-between gap-2">
          <span className="text-lg font-bold text-primary">
            {formatCurrency(product.price)}
          </span>
          <Button
            size="sm"
            onClick={handleAddToCart}
            disabled={outOfStock || adding}
            className="shrink-0"
          >
            <ShoppingCart className="h-4 w-4 mr-1.5" />
            {adding ? 'Adding…' : 'Add to Cart'}
          </Button>
        </div>
      </div>
    </div>
  )
}
