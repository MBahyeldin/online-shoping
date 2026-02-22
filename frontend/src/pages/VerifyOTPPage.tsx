import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useState, useEffect } from 'react'
import { useNavigate, useLocation, Link } from 'react-router-dom'
import { Cake, Loader2, RefreshCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { authService } from '@/services/auth'
import { useAuthStore } from '@/store/authStore'

const schema = z.object({
  email: z.string().email('Enter a valid email address'),
  otp: z.string().length(6, 'OTP must be exactly 6 digits').regex(/^\d+$/, 'OTP must be numeric'),
})

type FormValues = z.infer<typeof schema>

export function VerifyOTPPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const { setAuth } = useAuthStore()
  const [serverError, setServerError] = useState('')
  const [successMsg, setSuccessMsg] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [isResending, setIsResending] = useState(false)
  const [cooldown, setCooldown] = useState(0)

  const prefillEmail = (location.state as { email?: string })?.email ?? ''

  const {
    register,
    handleSubmit,
    getValues,
    setValue,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: prefillEmail },
  })

  // Cooldown timer for resend
  useEffect(() => {
    if (cooldown <= 0) return
    const timer = setTimeout(() => setCooldown((c) => c - 1), 1000)
    return () => clearTimeout(timer)
  }, [cooldown])

  const onSubmit = async (values: FormValues) => {
    setServerError('')
    setIsLoading(true)
    try {
      const result = await authService.verifyOTP(values)
      setAuth(result.user, result.token)
      navigate('/products', { replace: true })
    } catch (err: unknown) {
      setServerError(err instanceof Error ? err.message : 'Verification failed. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const handleResend = async () => {
    const email = getValues('email')
    if (!email) {
      setServerError('Please enter your email address first.')
      return
    }
    setIsResending(true)
    setServerError('')
    setSuccessMsg('')
    try {
      await authService.resendOTP(email)
      setSuccessMsg('A new OTP has been sent to your email.')
      setCooldown(60)
    } catch (err: unknown) {
      setServerError(err instanceof Error ? err.message : 'Failed to resend OTP.')
    } finally {
      setIsResending(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-orange-50 px-4">
      <div className="w-full max-w-md">
        {/* Card */}
        <div className="bg-white rounded-2xl shadow-sm border p-8">
          {/* Logo */}
          <Link to="/" className="flex items-center justify-center gap-2 text-primary font-bold text-xl mb-6">
            <Cake className="h-6 w-6" />
            Cake Shop
          </Link>

          <div className="text-center mb-8">
            <div className="text-4xl mb-3">📬</div>
            <h1 className="text-2xl font-extrabold text-gray-900">Check your email</h1>
            <p className="text-sm text-muted-foreground mt-2">
              We sent a 6-digit verification code to your email address. It expires in 5 minutes.
            </p>
          </div>

          {serverError && (
            <div className="mb-4 rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-sm text-destructive">
              {serverError}
            </div>
          )}

          {successMsg && (
            <div className="mb-4 rounded-lg bg-green-50 border border-green-200 p-3 text-sm text-green-700">
              {successMsg}
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5" noValidate>
            <div className="space-y-1.5">
              <Label htmlFor="email">Email Address</Label>
              <Input
                id="email"
                type="email"
                placeholder="you@example.com"
                {...register('email')}
                aria-invalid={!!errors.email}
              />
              {errors.email && (
                <p className="text-xs text-destructive">{errors.email.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="otp">Verification Code</Label>
              <Input
                id="otp"
                type="text"
                inputMode="numeric"
                maxLength={6}
                placeholder="123456"
                className="text-center text-2xl tracking-[0.5em] font-bold"
                {...register('otp')}
                aria-invalid={!!errors.otp}
              />
              {errors.otp && (
                <p className="text-xs text-destructive">{errors.otp.message}</p>
              )}
            </div>

            <Button type="submit" className="w-full" size="lg" disabled={isLoading}>
              {isLoading && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
              {isLoading ? 'Verifying…' : 'Verify & Sign In'}
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-muted-foreground">
              Didn't receive the code?{' '}
              <button
                onClick={handleResend}
                disabled={isResending || cooldown > 0}
                className="text-primary font-medium hover:underline disabled:opacity-50 disabled:no-underline inline-flex items-center gap-1"
              >
                {isResending && <RefreshCw className="h-3 w-3 animate-spin" />}
                {cooldown > 0 ? `Resend in ${cooldown}s` : 'Resend OTP'}
              </button>
            </p>
            <p className="mt-3 text-sm">
              <Link to="/register" className="text-primary font-medium hover:underline">
                ← Back to Registration
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
