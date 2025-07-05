
-- name: CreateProject :one
INSERT INTO projects ( user_id, name, description, color) VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetProjectById :one
SELECT * FROM projects WHERE id=$1 ;

-- name: GetProjectsByUserId :many
SELECT * FROM projects WHERE user_id=$1 ORDER BY id;

-- name: GetProjectByName :one

SELECT * FROM projects WHERE name=$1 AND user_id=$2;

-- name: UpdateProject :one
UPDATE projects SET name=$2, description=$3, color=$4, updated_at=now()
WHERE id=$1 RETURNING *;    

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;


