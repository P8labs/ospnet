import { useMutation } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"
import type { OnboardingTokenResponse } from "@/types/onboarding"

export function useCreateToken() {
  return useMutation({
    mutationFn: async (): Promise<OnboardingTokenResponse> => {
      return apiClient.createOnboardingToken()
    },
  })
}
