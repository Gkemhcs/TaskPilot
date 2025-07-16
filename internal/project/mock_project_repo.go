package project

import (
	"context"

	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	"github.com/stretchr/testify/mock"
)

// MockProjectRepo is a mock implementation of the projectdb.Querier interface
type MockProjectRepo struct {
	mock.Mock
}

func (m *MockProjectRepo) CreateProject(ctx context.Context, arg projectdb.CreateProjectParams) (projectdb.Project, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(projectdb.Project), args.Error(1)
}

func (m *MockProjectRepo) DeleteProject(ctx context.Context, projectID int64) error {
	args := m.Called(ctx, projectID)
	return args.Error(0)
}

func (m *MockProjectRepo) GetProjectById(ctx context.Context, id int64) (projectdb.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(projectdb.Project), args.Error(1)
}

func (m *MockProjectRepo) GetProjectsByUserId(ctx context.Context, userID int32) ([]projectdb.Project, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]projectdb.Project), args.Error(1)
}

func (m *MockProjectRepo) GetProjectByName(ctx context.Context, params projectdb.GetProjectByNameParams) (projectdb.Project, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(projectdb.Project), args.Error(1)
}

func (m *MockProjectRepo) UpdateProject(ctx context.Context, arg projectdb.UpdateProjectParams)  error{
	args:=m.Called(ctx,arg)
	return args.Error(0)
}

