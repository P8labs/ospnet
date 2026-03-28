-- name: CreateToken :exec
INSERT INTO onboarding_tokens (
    token,
    expires_at
) VALUES (?, ?);


-- name: GetToken :one
SELECT *
FROM onboarding_tokens
WHERE token = ?;


-- name: MarkTokenUsed :exec
UPDATE onboarding_tokens
SET used = TRUE
WHERE token = ?;


-- name: DeleteExpiredTokens :exec
DELETE FROM onboarding_tokens
WHERE expires_at < CURRENT_TIMESTAMP;


-- name: GetValidToken :one
SELECT *
FROM onboarding_tokens
WHERE token = ?
  AND used = FALSE
  AND expires_at > CURRENT_TIMESTAMP;