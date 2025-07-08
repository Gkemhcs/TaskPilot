
-- name: CreateTask :one

INSERT INTO tasks (project_id, assignee_id, title, description, status, priority, due_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTaskById :one
SELECT * FROM tasks WHERE id = $1;

-- name: GetTasksByProjectId :many
SELECT * FROM tasks WHERE project_id = $1 ORDER BY id;




-- name: GetAllTasks :many
SELECT * FROM tasks ORDER BY id;

-- name: DeleteTask :execrows
DELETE FROM tasks WHERE id = $1;




-- name: UpdateTask :execrows

UPDATE tasks
SET
  title = COALESCE(sqlc.narg('title'), title),
  description = COALESCE(sqlc.narg('description'), description),
  due_date = COALESCE(sqlc.narg('due_date'), due_date),
  status = COALESCE(sqlc.narg('status'), status),
  priority = COALESCE(sqlc.narg('priority'), priority),
  updated_at = now()
WHERE id = sqlc.arg('id');