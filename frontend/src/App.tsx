import { Routes, Route } from 'react-router-dom'
import { Header } from '@/components/layout/Header'
import { Footer } from '@/components/layout/Footer'
import { CartDrawer } from '@/components/shared/CartDrawer'
import { HomePage } from '@/pages/HomePage'
import { RegisterPage } from '@/pages/RegisterPage'
import { VerifyOTPPage } from '@/pages/VerifyOTPPage'
import { ProductsPage } from '@/pages/ProductsPage'
import { CartPage } from '@/pages/CartPage'

function App() {
  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <CartDrawer />

      <div className="flex-1">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/verify-otp" element={<VerifyOTPPage />} />
          <Route path="/products" element={<ProductsPage />} />
          <Route path="/cart" element={<CartPage />} />

          {/* 404 */}
          <Route
            path="*"
            element={
              <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4 text-center px-4">
                <p className="text-6xl font-extrabold text-primary">404</p>
                <h1 className="text-2xl font-bold text-gray-900">Page Not Found</h1>
                <p className="text-muted-foreground">
                  The page you're looking for doesn't exist.
                </p>
                <a href="/" className="text-primary underline underline-offset-4">
                  Go back home
                </a>
              </div>
            }
          />
        </Routes>
      </div>

      <Footer />
    </div>
  )
}

export default App
