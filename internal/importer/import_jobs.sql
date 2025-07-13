-- name: CreateImportJob :one
INSERT INTO import_jobs (
  id, file_path, importer_type, user_id, status
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateImportJobStatus :exec
UPDATE import_jobs
SET status = $3,
    error_message = $4,
    updated_at = NOW()
WHERE id = $1 and user_id=$2;

-- name: GetImportJob :one
SELECT * FROM import_jobs
WHERE id = $1 and user_id = $2;

-- name: ListImportJobs :many
SELECT * FROM import_jobs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
