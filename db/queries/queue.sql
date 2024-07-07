-- name: FindQueue :many
SELECT * FROM "queue";

-- name: FindQueueByID :one
SELECT * FROM "queue" WHERE id=$1 AND user_id=$2
LIMIT 1 FOR UPDATE;

-- name: FindLatestQueue :one
SELECT * FROM "queue"
ORDER BY created_at DESC
LIMIT 1 FOR UPDATE NOWAIT;

-- name: UpdateQueue :one
UPDATE "queue"
SET service_time=$2,
total_waiting_time=$3
WHERE id=$1 RETURNING *;

-- name: CreateQueue :one
INSERT INTO "queue"(queue_number, user_id, arrival_time) VALUES
($1, $2, $3) RETURNING *;

-- name: DeleteQueue :exec
DELETE FROM "queue" WHERE id=$1 AND user_id=$2;
