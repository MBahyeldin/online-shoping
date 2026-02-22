-- name: CreateOTP :one
INSERT INTO email_otps (user_id, otp_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetLatestOTPByUserID :one
SELECT * FROM email_otps
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkOTPUsed :exec
UPDATE email_otps
SET is_used = TRUE
WHERE id = $1;

-- name: IncrementOTPAttempts :one
UPDATE email_otps
SET attempt_count = attempt_count + 1
WHERE id = $1
RETURNING *;

-- name: CountRecentOTPsByUserID :one
SELECT COUNT(*) FROM email_otps
WHERE user_id = $1
  AND created_at > NOW() - INTERVAL '1 hour';

-- name: InvalidateUserOTPs :exec
UPDATE email_otps
SET is_used = TRUE
WHERE user_id = $1 AND is_used = FALSE;
