
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





-- name: ListTasksWithFilters :many
SELECT *
FROM tasks
WHERE 
    (project_id = COALESCE(sqlc.narg('project_id'), project_id))
  AND (assignee_id = COALESCE(sqlc.narg('assignee_id'), assignee_id))
  AND (status = ANY(sqlc.narg('statuses')))
  AND (priority = COALESCE(sqlc.narg('priority'), priority))
  AND (due_date >= COALESCE(sqlc.narg('due_date_from'), due_date))
  AND (due_date <= COALESCE(sqlc.narg('due_date_to'), due_date))
ORDER BY due_date
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
