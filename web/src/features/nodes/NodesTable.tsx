import { useNodesData } from "@/hooks"
import type { Node } from "@/types/node"
import { StatusBadge } from "@/components/StatusBadge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Skeleton } from "@/components/ui/skeleton"

const formatDate = (dateString: string) => {
  try {
    return new Date(dateString).toLocaleString()
  } catch {
    return dateString
  }
}

export function NodesTable() {
  const { nodes, isLoading, isError, error } = useNodesData()

  if (isLoading) {
    return (
      <div className="rounded-lg border border-border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Hostname</TableHead>
              <TableHead>Region</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>CPU / Memory</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Last Seen</TableHead>
              <TableHead>IP Address</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {[...Array(5)].map((_, i) => (
              <TableRow key={i}>
                {[...Array(8)].map((_, j) => (
                  <TableCell key={j}>
                    <Skeleton className="h-4 w-24" />
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    )
  }

  if (isError) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-800">
        <p className="font-semibold">Failed to load nodes</p>
        <p className="mt-1 text-xs">
          {error instanceof Error ? error.message : "Unknown error occurred"}
        </p>
      </div>
    )
  }

  if (nodes.length === 0) {
    return (
      <div className="flex h-64 flex-col items-center justify-center rounded-lg border border-dashed border-border bg-background p-8">
        <div className="text-center">
          <h3 className="font-semibold text-foreground">No nodes yet</h3>
          <p className="mt-1 text-sm text-muted-foreground">
            Add your first node to get started
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow className="bg-muted/50 hover:bg-muted/50">
            <TableHead className="font-semibold">Name</TableHead>
            <TableHead className="font-semibold">Hostname</TableHead>
            <TableHead className="font-semibold">Region</TableHead>
            <TableHead className="font-semibold">Type</TableHead>
            <TableHead className="font-semibold">CPU / Memory</TableHead>
            <TableHead className="font-semibold">Status</TableHead>
            <TableHead className="font-semibold">Last Seen</TableHead>
            <TableHead className="font-semibold">IP Address</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {nodes.map((node: Node) => (
            <TableRow key={node.id} className="hover:bg-muted/50">
              <TableCell className="font-medium">{node.name}</TableCell>
              <TableCell className="text-xs text-muted-foreground">
                {node.hostname}
              </TableCell>
              <TableCell className="text-sm">{node.region}</TableCell>
              <TableCell className="text-sm">{node.type}</TableCell>
              <TableCell className="text-sm">
                {node.cpu} CPU / {(node.memory / 1024 / 1024 / 1024).toFixed(1)}{" "}
                GB
              </TableCell>
              <TableCell>
                <StatusBadge status={node.status} />
              </TableCell>
              <TableCell className="text-xs text-muted-foreground">
                {formatDate(node.last_seen)}
              </TableCell>
              <TableCell className="font-mono text-xs">{node.ip}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
