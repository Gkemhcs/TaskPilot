
-- name: CreateProject :one
INSERT INTO projects ( user_id, name, description, color) VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetProjectById :one
SELECT * FROM projects WHERE id=$1 ;

-- name: GetProjectsByUserId :many
SELECT * FROM projects WHERE user_id=$1 ORDER BY id;

-- name: GetProjectByName :one

SELECT * FROM projects WHERE name=$1 AND user_id=$2;


-- name: UpdateProject :exec
UPDATE projects
SET
  name = COALESCE(sqlc.narg('name'), name),
  description = COALESCE(sqlc.narg('description'), description),
  color = COALESCE(sqlc.narg('color'), color) ,
  updated_at = now()
WHERE id = sqlc.arg('id');


-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;


