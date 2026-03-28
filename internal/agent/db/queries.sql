-- name: CreateContainer :exec
INSERT INTO containers (id, image, name, port, status)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateContainerStatus :exec
UPDATE containers
SET status = ?, docker_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: GetContainers :many
SELECT * FROM containers;

-- name: GetContainerByID :one
SELECT * FROM containers WHERE id = ?;