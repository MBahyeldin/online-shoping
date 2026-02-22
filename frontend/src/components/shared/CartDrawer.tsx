import { X, Plus, Minus, Trash2, ShoppingBag } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { useCartStore } from '@/store/cartStore'
import { cartService } from '@/services/cart'
import { formatCurrency, cn } from '@/lib/utils'

export function CartDrawer() {
  const { cart, isOpen, closeCart, setCart } = useCartStore()
  const [loadingItem, setLoadingItem] = useState<string | null>(null)

  const handleUpdateQuantity = async (itemId: string, newQty: number) => {
    setLoadingItem(itemId)
    try {
      if (newQty < 1) {
        const updated = await cartService.removeItem(itemId)
        setCart(updated)
      } else {
        const updated = await cartService.updateItem(itemId, newQty)
        setCart(updated)
      }
    } catch (err) {
      console.error('Update cart error:', err)
    } finally {
      setLoadingItem(null)
    }
  }

  const handleRemove = async (itemId: string) => {
    setLoadingItem(itemId)
    try {
      const updated = await cartService.removeItem(itemId)
      setCart(updated)
    } catch (err) {
      console.error('Remove item error:', err)
    } finally {
      setLoadingItem(null)
    }
  }

  return (
    <>
      {/* Backdrop */}
      <div
        className={cn(
          'fixed inset-0 z-40 bg-black/50 transition-opacity duration-300',
          isOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'
        )}
        onClick={closeCart}
      />

      {/* Drawer */}
      <aside
        className={cn(
          'fixed right-0 top-0 z-50 h-full w-full max-w-sm bg-white shadow-2xl flex flex-col transition-transform duration-300',
          isOpen ? 'translate-x-0' : 'translate-x-full'
        )}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b">
          <h2 className="text-lg font-bold">My Cart</h2>
          <Button variant="ghost" size="icon" onClick={closeCart}>
            <X className="h-5 w-5" />
          </Button>
        </div>

        {/* Items */}
        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {!cart || cart.items.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full gap-4 text-center">
              <ShoppingBag className="h-16 w-16 text-muted-foreground/20" />
              <p className="text-muted-foreground font-medium">Your cart is empty</p>
              <Button variant="outline" onClick={closeCart} asChild>
                <Link to="/products">Browse Cakes</Link>
              </Button>
            </div>
          ) : (
            cart.items.map((item) => (
              <div key={item.id} className="flex gap-3">
                {/* Product image */}
                <div className="h-16 w-16 shrink-0 rounded-lg overflow-hidden bg-muted">
                  {item.product_image_url ? (
                    <img
                      src={item.product_image_url}
                      alt={item.product_name}
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <div className="h-full w-full flex items-center justify-center text-2xl">
                      🎂
                    </div>
                  )}
                </div>

                {/* Details */}
                <div className="flex-1 min-w-0">
                  <p className="font-medium text-sm truncate">{item.product_name}</p>
                  <p className="text-primary font-bold text-sm">{formatCurrency(item.price)}</p>

                  {/* Quantity controls */}
                  <div className="flex items-center gap-2 mt-1.5">
                    <button
                      onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                      disabled={loadingItem === item.id}
                      className="h-6 w-6 flex items-center justify-center rounded border hover:bg-muted transition-colors disabled:opacity-50"
                    >
                      <Minus className="h-3 w-3" />
                    </button>
                    <span className="text-sm font-semibold w-6 text-center">{item.quantity}</span>
                    <button
                      onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                      disabled={loadingItem === item.id}
                      className="h-6 w-6 flex items-center justify-center rounded border hover:bg-muted transition-colors disabled:opacity-50"
                    >
                      <Plus className="h-3 w-3" />
                    </button>
                  </div>
                </div>

                {/* Subtotal + delete */}
                <div className="flex flex-col items-end justify-between">
                  <button
                    onClick={() => handleRemove(item.id)}
                    disabled={loadingItem === item.id}
                    className="text-muted-foreground hover:text-destructive transition-colors disabled:opacity-50"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                  <p className="text-sm font-semibold">{formatCurrency(item.subtotal)}</p>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Footer */}
        {cart && cart.items.length > 0 && (
          <div className="border-t p-4 space-y-3">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Subtotal</span>
              <span className="font-semibold">{formatCurrency(cart.total)}</span>
            </div>
            <Separator />
            <div className="flex justify-between font-bold text-base">
              <span>Total</span>
              <span className="text-primary">{formatCurrency(cart.total)}</span>
            </div>
            <Button className="w-full" size="lg" asChild onClick={closeCart}>
              <Link to="/cart">Proceed to Checkout</Link>
            </Button>
          </div>
        )}
      </aside>
    </>
  )
}
