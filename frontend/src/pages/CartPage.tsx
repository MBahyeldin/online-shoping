import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link, useNavigate } from 'react-router-dom'
import { ShoppingBag, Minus, Plus, Trash2, Loader2, CheckCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import { useCartStore } from '@/store/cartStore'
import { useAuthStore } from '@/store/authStore'
import { cartService } from '@/services/cart'
import { orderService } from '@/services/orders'
import { formatCurrency, getMinDeliveryDate } from '@/lib/utils'
import { PageLoader } from '@/components/shared/LoadingSpinner'

const checkoutSchema = z.object({
  delivery_address: z.string().min(10, 'Please enter a full delivery address'),
  delivery_date: z.string().min(1, 'Delivery date is required'),
  notes: z.string().optional(),
  payment_method: z.literal('cash_on_delivery'),
})

type CheckoutFormValues = z.infer<typeof checkoutSchema>

export function CartPage() {
  const { cart, setCart, clearCart } = useCartStore()
  const { isAuthenticated } = useAuthStore()
  const navigate = useNavigate()

  const [isLoadingCart, setIsLoadingCart] = useState(true)
  const [loadingItem, setLoadingItem] = useState<string | null>(null)
  const [isOrdering, setIsOrdering] = useState(false)
  const [orderSuccess, setOrderSuccess] = useState<string | null>(null)
  const [orderError, setOrderError] = useState('')

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CheckoutFormValues>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: { payment_method: 'cash_on_delivery' },
  })

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/register')
      return
    }
    cartService
      .getCart()
      .then(setCart)
      .catch(console.error)
      .finally(() => setIsLoadingCart(false))
  }, [isAuthenticated, navigate, setCart])

  const updateQty = async (itemId: string, qty: number) => {
    setLoadingItem(itemId)
    try {
      if (qty < 1) {
        const updated = await cartService.removeItem(itemId)
        setCart(updated)
      } else {
        const updated = await cartService.updateItem(itemId, qty)
        setCart(updated)
      }
    } finally {
      setLoadingItem(null)
    }
  }

  const removeItem = async (itemId: string) => {
    setLoadingItem(itemId)
    try {
      const updated = await cartService.removeItem(itemId)
      setCart(updated)
    } finally {
      setLoadingItem(null)
    }
  }

  const onCheckout = async (values: CheckoutFormValues) => {
    if (!cart || cart.items.length === 0) return
    setIsOrdering(true)
    setOrderError('')
    try {
      // Convert local datetime to RFC3339
      const localDate = new Date(values.delivery_date)
      const deliveryISO = localDate.toISOString()

      const order = await orderService.create({
        delivery_address: values.delivery_address,
        delivery_date: deliveryISO,
        notes: values.notes || undefined,
        payment_method: values.payment_method,
      })
      clearCart()
      setOrderSuccess(order.id)
    } catch (err: unknown) {
      setOrderError(err instanceof Error ? err.message : 'Failed to place order.')
    } finally {
      setIsOrdering(false)
    }
  }

  if (!isAuthenticated) return null

  if (isLoadingCart) return <PageLoader />

  // Order success screen
  if (orderSuccess) {
    return (
      <main className="min-h-screen flex items-center justify-center bg-orange-50 px-4">
        <div className="max-w-md w-full text-center space-y-6 bg-white rounded-2xl p-10 shadow-sm border">
          <CheckCircle className="h-16 w-16 text-green-500 mx-auto" />
          <div>
            <h1 className="text-2xl font-extrabold text-gray-900">Order Placed!</h1>
            <p className="text-muted-foreground mt-2">
              Your order has been placed successfully. We'll prepare it with care!
            </p>
            <p className="text-xs text-muted-foreground mt-1">Order ID: {orderSuccess}</p>
          </div>
          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <Button asChild>
              <Link to="/products">Continue Shopping</Link>
            </Button>
          </div>
        </div>
      </main>
    )
  }

  const isEmpty = !cart || cart.items.length === 0

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="container mx-auto max-w-6xl px-4 py-10">
        <h1 className="text-3xl font-extrabold text-gray-900 mb-8">My Cart</h1>

        {isEmpty ? (
          <div className="flex flex-col items-center justify-center py-24 gap-5 text-center">
            <ShoppingBag className="h-20 w-20 text-muted-foreground/20" />
            <p className="text-xl font-semibold text-gray-900">Your cart is empty</p>
            <p className="text-muted-foreground">Add some delicious cakes to get started!</p>
            <Button asChild>
              <Link to="/products">Browse Cakes</Link>
            </Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Cart items */}
            <div className="lg:col-span-2 space-y-4">
              {cart.items.map((item) => (
                <div
                  key={item.id}
                  className="flex gap-4 bg-white rounded-xl border p-4 shadow-sm"
                >
                  {/* Image */}
                  <div className="h-20 w-20 shrink-0 rounded-lg overflow-hidden bg-muted">
                    {item.product_image_url ? (
                      <img
                        src={item.product_image_url}
                        alt={item.product_name}
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <div className="h-full w-full flex items-center justify-center text-3xl">🎂</div>
                    )}
                  </div>

                  {/* Details */}
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold truncate">{item.product_name}</p>
                    <p className="text-sm text-primary font-bold">{formatCurrency(item.price)} / each</p>

                    {/* Quantity controls */}
                    <div className="flex items-center gap-2 mt-2">
                      <button
                        onClick={() => updateQty(item.id, item.quantity - 1)}
                        disabled={loadingItem === item.id}
                        className="h-7 w-7 flex items-center justify-center rounded border hover:bg-muted disabled:opacity-50"
                      >
                        <Minus className="h-3 w-3" />
                      </button>
                      <span className="text-sm font-bold w-6 text-center">{item.quantity}</span>
                      <button
                        onClick={() => updateQty(item.id, item.quantity + 1)}
                        disabled={loadingItem === item.id}
                        className="h-7 w-7 flex items-center justify-center rounded border hover:bg-muted disabled:opacity-50"
                      >
                        <Plus className="h-3 w-3" />
                      </button>
                    </div>
                  </div>

                  {/* Subtotal + remove */}
                  <div className="flex flex-col items-end justify-between">
                    <button
                      onClick={() => removeItem(item.id)}
                      disabled={loadingItem === item.id}
                      className="text-muted-foreground hover:text-destructive transition-colors disabled:opacity-50"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                    <p className="font-bold">{formatCurrency(item.subtotal)}</p>
                  </div>
                </div>
              ))}
            </div>

            {/* Order summary + checkout form */}
            <div className="space-y-6">
              {/* Summary */}
              <div className="bg-white rounded-xl border p-5 shadow-sm space-y-4">
                <h2 className="font-bold text-lg">Order Summary</h2>
                <div className="space-y-2 text-sm">
                  {cart.items.map((item) => (
                    <div key={item.id} className="flex justify-between text-muted-foreground">
                      <span className="truncate max-w-[200px]">
                        {item.product_name} × {item.quantity}
                      </span>
                      <span className="font-medium text-foreground ml-2">
                        {formatCurrency(item.subtotal)}
                      </span>
                    </div>
                  ))}
                </div>
                <Separator />
                <div className="flex justify-between font-bold text-base">
                  <span>Total</span>
                  <span className="text-primary">{formatCurrency(cart.total)}</span>
                </div>
              </div>

              {/* Checkout form */}
              <div className="bg-white rounded-xl border p-5 shadow-sm">
                <h2 className="font-bold text-lg mb-5">Delivery Details</h2>

                {orderError && (
                  <div className="mb-4 rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-sm text-destructive">
                    {orderError}
                  </div>
                )}

                <form onSubmit={handleSubmit(onCheckout)} className="space-y-4">
                  <div className="space-y-1.5">
                    <Label htmlFor="delivery_address">Delivery Address</Label>
                    <Textarea
                      id="delivery_address"
                      placeholder="123 Main St, City, State, ZIP"
                      rows={3}
                      {...register('delivery_address')}
                      aria-invalid={!!errors.delivery_address}
                    />
                    {errors.delivery_address && (
                      <p className="text-xs text-destructive">{errors.delivery_address.message}</p>
                    )}
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="delivery_date">Delivery Date & Time</Label>
                    <Input
                      id="delivery_date"
                      type="datetime-local"
                      min={getMinDeliveryDate()}
                      {...register('delivery_date')}
                      aria-invalid={!!errors.delivery_date}
                    />
                    {errors.delivery_date && (
                      <p className="text-xs text-destructive">{errors.delivery_date.message}</p>
                    )}
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="notes">Custom Message / Notes (optional)</Label>
                    <Textarea
                      id="notes"
                      placeholder="e.g. 'Happy Birthday Sarah!' or special instructions"
                      rows={2}
                      {...register('notes')}
                    />
                  </div>

                  {/* Payment method */}
                  <div className="rounded-lg border bg-muted/40 p-3 flex items-center gap-3">
                    <div className="h-8 w-8 rounded-full bg-green-100 flex items-center justify-center text-base">
                      💵
                    </div>
                    <div>
                      <p className="text-sm font-semibold">Cash on Delivery</p>
                      <p className="text-xs text-muted-foreground">Pay when your order arrives</p>
                    </div>
                    <input type="hidden" {...register('payment_method')} />
                  </div>

                  <Button
                    type="submit"
                    className="w-full"
                    size="lg"
                    disabled={isOrdering}
                  >
                    {isOrdering && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
                    {isOrdering ? 'Placing Order…' : `Place Order · ${formatCurrency(cart.total)}`}
                  </Button>
                </form>
              </div>
            </div>
          </div>
        )}
      </div>
    </main>
  )
}
