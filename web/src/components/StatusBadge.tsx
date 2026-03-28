import type { NodeStatus } from "@/types/node"
import { Badge } from "@/components/ui/badge"

interface StatusBadgeProps {
  status: NodeStatus
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const config = {
    healthy: {
      label: "Healthy",
      variant: "default" as const,
      className: "bg-green-100 text-green-800 hover:bg-green-100",
    },
    unhealthy: {
      label: "Unhealthy",
      variant: "default" as const,
      className: "bg-yellow-100 text-yellow-800 hover:bg-yellow-100",
    },
    offline: {
      label: "Offline",
      variant: "default" as const,
      className: "bg-red-100 text-red-800 hover:bg-red-100",
    },
    unknown: {
      label: "Unknown",
      variant: "default" as const,
      className: "bg-gray-100 text-gray-800 hover:bg-gray-100",
    },
  }

  const { label, className } = config[status]

  return <Badge className={className}>{label}</Badge>
}
