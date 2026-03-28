import { useQuery } from "@tanstack/react-query"
import { apiClient } from "@/lib/api"
import type { Node } from "@/types/node"

export function useNodes() {
  return useQuery({
    queryKey: ["nodes"],
    queryFn: async () => {
      const response = await apiClient.getNodes()
      return response || []
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
  })
}

export function useNodesData() {
  const { data = [], isLoading, error, isError } = useNodes()

  return {
    nodes: data as Node[],
    isLoading,
    error,
    isError,
  }
}
