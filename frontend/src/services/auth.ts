import api from '@/lib/api'
import type { AuthResult, RegisterPayload, VerifyOTPPayload } from '@/types'

export const authService = {
  register: async (payload: RegisterPayload): Promise<void> => {
    await api.post('/auth/register', payload)
  },

  verifyOTP: async (payload: VerifyOTPPayload): Promise<AuthResult> => {
    const { data } = await api.post<{ success: boolean; data: AuthResult }>(
      '/auth/verify-otp',
      payload
    )
    return data.data!
  },

  resendOTP: async (email: string): Promise<void> => {
    await api.post('/auth/resend-otp', { email })
  },
}
