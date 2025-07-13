-- name: CreateExportJob :one
INSERT INTO export_jobs (id, user_id, export_type)
VALUES ($1, $2, $3)
RETURNING id, user_id, export_type, status, url, error_message, created_at, updated_at;

-- name: UpdateExportJobStatus :exec
UPDATE export_jobs
SET status = $2, updated_at = NOW(), error_message = $3
WHERE id = $1 AND user_id = $4;

-- name: UpdateExportJobURL :exec
UPDATE export_jobs
SET status = 'completed', url = $2, updated_at = NOW()
WHERE id = $1 AND user_id = $3;


-- name: GetExportJobStatus :one


SELECT * FROM export_jobs 
WHERE user_id=$1 AND id=$2 ;
