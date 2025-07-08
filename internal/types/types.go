package types

import (
	"context"

	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
)

type IUserHandler interface {
	CreateUser()
	LoginUser()
	CheckHashedPassword()
}

type TaskQueryService interface {
	GetTasksByProjectID(ctx context.Context, projectID int) ([]taskdb.Task, error)
}

type ProjectReader interface {
	GetProjectById(ctx context.Context, projectId int) (*projectdb.Project, error)
}
