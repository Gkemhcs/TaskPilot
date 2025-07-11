package task

import (
	"context"

	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) CreateTask(ctx context.Context, arg taskdb.CreateTaskParams) (taskdb.Task, error) {

	args := m.Called(ctx, arg)
	return args.Get(0).(taskdb.Task), args.Error(1)

}

func (m *MockTaskRepo) GetTaskById(ctx context.Context, arg int64) (taskdb.Task, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(taskdb.Task), args.Error(1)

}
func (m *MockTaskRepo) GetTasksByProjectId(ctx context.Context, arg int64) ([]taskdb.Task, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]taskdb.Task), args.Error(1)

}

func (m *MockTaskRepo) DeleteTask(ctx context.Context, arg int64) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func(m *MockTaskRepo)UpdateTask(ctx context.Context, arg taskdb.UpdateTaskParams) (int64, error){
	args:=m.Called(ctx,arg)
	return args.Get(0).(int64),args.Error(1)
}
func(m *MockTaskRepo) ListTasksWithFilters(ctx context.Context, arg taskdb.ListTasksWithFiltersParams) ([]taskdb.Task, error){
	args:=m.Called(ctx,arg)
	return args.Get(0).([]taskdb.Task),args.Error(1)
}
func(m *MockTaskRepo) GetAllTasks(ctx context.Context) ([]taskdb.Task, error){
	args:=m.Called(ctx)
	return args.Get(0).([]taskdb.Task),args.Error(1)
}

