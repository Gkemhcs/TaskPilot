package user

import (
	"context"

	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}



func (m *MockUserRepo) CreateUser(ctx context.Context, arg userdb.CreateUserParams) (userdb.User, error){
	args:=m.Called(ctx,arg)
	return args.Get(0).(userdb.User),args.Error(1)


}

func(m *MockUserRepo) DeleteUser(ctx context.Context, id int32) error{
	args:=m.Called(ctx,id)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (userdb.User, error){
	args:=m.Called(ctx,email)
	return args.Get(0).(userdb.User),args.Error(1)
}
func(m *MockUserRepo)GetUserById(ctx context.Context, id int32) (userdb.User, error){
	args:=m.Called(ctx,id)
	return args.Get(0).(userdb.User),args.Error(1)
}

func(m *MockUserRepo)	GetUserByName(ctx context.Context, name string) (userdb.User, error){
	args:=m.Called(ctx,name)
	return args.Get(0).(userdb.User),args.Error(1)
}

func(m *MockUserRepo)ListUsers(ctx context.Context) ([]userdb.User, error){
args:=m.Called(ctx)
return args.Get(0).([]userdb.User),args.Error(1)
}