import type {
  GetNodesResponse,
  RegisterNodeRequest,
  RegisterNodeResponse,
} from "@/types/node"
import type {
  OnboardingTokenRequest,
  OnboardingTokenResponse,
} from "@/types/onboarding"

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8000/api"

class APIClient {
  private baseURL: string

  constructor(baseURL: string) {
    this.baseURL = baseURL
  }

  private async fetch<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`

    const response = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({}))
      throw new Error(
        error.message || `API Error: ${response.status} ${response.statusText}`
      )
    }

    const data = await response.json()

    if (data.status != true) {
      throw Error(data.message || "Problem in response")
    }

    return data.data as Promise<T>
  }

  async createOnboardingToken(
    _payload?: OnboardingTokenRequest
  ): Promise<OnboardingTokenResponse> {
    return this.fetch<OnboardingTokenResponse>("/onboarding/token", {
      method: "POST",
      body: JSON.stringify(_payload || {}),
    })
  }

  async registerNode(
    payload: RegisterNodeRequest
  ): Promise<RegisterNodeResponse> {
    return this.fetch<RegisterNodeResponse>("/onboarding/register", {
      method: "POST",
      body: JSON.stringify(payload),
    })
  }

  async getNodes(): Promise<GetNodesResponse> {
    return this.fetch<GetNodesResponse>("/nodes", {
      method: "GET",
    })
  }
}

export const apiClient = new APIClient(API_BASE_URL)
