import { Link, useNavigate } from 'react-router-dom'
import { ShoppingCart, Cake, User, LogOut, Menu, X } from 'lucide-react'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/authStore'
import { useCartStore } from '@/store/cartStore'
import { cn } from '@/lib/utils'

export function Header() {
  const { isAuthenticated, user, clearAuth } = useAuthStore()
  const { itemCount, openCart } = useCartStore()
  const navigate = useNavigate()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)

  const count = itemCount()

  const handleLogout = () => {
    clearAuth()
    navigate('/')
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-white/95 backdrop-blur supports-[backdrop-filter]:bg-white/60">
      <div className="container mx-auto flex h-16 max-w-6xl items-center justify-between px-4">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-2 font-bold text-xl text-primary">
          <Cake className="h-6 w-6" />
          <span>Cake Shop</span>
        </Link>

        {/* Desktop Nav */}
        <nav className="hidden md:flex items-center gap-6">
          <Link
            to="/"
            className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
          >
            Home
          </Link>
          <Link
            to="/products"
            className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
          >
            Our Cakes
          </Link>
        </nav>

        {/* Actions */}
        <div className="flex items-center gap-2">
          {/* Cart button (only when authenticated) */}
          {isAuthenticated && (
            <Button
              variant="ghost"
              size="icon"
              onClick={openCart}
              className="relative"
              aria-label="Shopping cart"
            >
              <ShoppingCart className="h-5 w-5" />
              {count > 0 && (
                <span className="absolute -top-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-primary-foreground">
                  {count > 99 ? '99+' : count}
                </span>
              )}
            </Button>
          )}

          {/* Auth buttons */}
          {isAuthenticated ? (
            <div className="hidden md:flex items-center gap-2">
              <span className="text-sm text-muted-foreground">
                Hi, {user?.first_name}
              </span>
              <Button variant="ghost" size="sm" onClick={handleLogout}>
                <LogOut className="h-4 w-4 mr-1" />
                Logout
              </Button>
            </div>
          ) : (
            <Button size="sm" asChild className="hidden md:inline-flex">
              <Link to="/register">
                <User className="h-4 w-4 mr-1" />
                Sign Up
              </Link>
            </Button>
          )}

          {/* Mobile menu toggle */}
          <Button
            variant="ghost"
            size="icon"
            className="md:hidden"
            onClick={() => setMobileMenuOpen((o) => !o)}
          >
            {mobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </Button>
        </div>
      </div>

      {/* Mobile menu */}
      <div
        className={cn(
          'md:hidden border-t bg-white transition-all duration-200 ease-in-out',
          mobileMenuOpen ? 'max-h-64 opacity-100' : 'max-h-0 opacity-0 overflow-hidden'
        )}
      >
        <nav className="container mx-auto flex flex-col gap-1 px-4 py-3">
          <Link
            to="/"
            onClick={() => setMobileMenuOpen(false)}
            className="py-2 text-sm font-medium hover:text-primary transition-colors"
          >
            Home
          </Link>
          <Link
            to="/products"
            onClick={() => setMobileMenuOpen(false)}
            className="py-2 text-sm font-medium hover:text-primary transition-colors"
          >
            Our Cakes
          </Link>
          {isAuthenticated ? (
            <button
              onClick={() => { handleLogout(); setMobileMenuOpen(false) }}
              className="py-2 text-left text-sm font-medium text-destructive hover:opacity-80 transition-opacity"
            >
              Logout
            </button>
          ) : (
            <Link
              to="/register"
              onClick={() => setMobileMenuOpen(false)}
              className="py-2 text-sm font-medium text-primary"
            >
              Sign Up / Login
            </Link>
          )}
        </nav>
      </div>
    </header>
  )
}
