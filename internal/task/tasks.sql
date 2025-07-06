
-- name: CreateTask :one

INSERT INTO tasks (project_id, assignee_id, title, description, status, priority, due_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTaskById :one
SELECT * FROM tasks WHERE id = $1;

-- name: GetTasksByProjectId :many
SELECT * FROM tasks WHERE project_id = $1 ORDER BY id;


-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;


