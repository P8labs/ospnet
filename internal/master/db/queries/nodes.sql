-- name: CreateNode :exec
INSERT INTO nodes (
    id,
    name,
    hostname,
    ip,
    cpu,
    memory,
    arch,
    region,
    type,
    status,
    last_seen
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);


-- name: GetNodeByID :one
SELECT *
FROM nodes
WHERE id = ?;


-- name: GetNodes :many
SELECT *
FROM nodes;


-- name: UpdateHeartbeat :exec
UPDATE nodes
SET last_seen = CURRENT_TIMESTAMP,
    status = 'healthy',
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;


-- name: MarkNodeUnhealthy :exec
UPDATE nodes
SET status = 'unhealthy',
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;


-- name: UpdateNodeIP :exec
UPDATE nodes
SET ip = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;


-- name: DeleteNode :exec
DELETE FROM nodes
WHERE id = ?;