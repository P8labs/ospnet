export type NodeStatus = "healthy" | "unhealthy" | "offline" | "unknown"

export interface Node {
  id: string
  name: string
  hostname: string
  region: string
  type: string
  cpu: number
  memory: number
  status: NodeStatus
  last_seen: string
  ip: string
}

export interface RegisterNodeRequest {
  token: string
  hostname: string
  ip: string
  region: string
  type: string
  cpu: number
  memory: number
}

export interface RegisterNodeResponse {
  id: string
  name: string
  status: NodeStatus
}

export type GetNodesResponse = Node[] | null
