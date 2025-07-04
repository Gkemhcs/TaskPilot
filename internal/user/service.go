package user

import (
	"context"
	"errors"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(dbQuerier userdb.Querier) *UserService {
	return &UserService{
		userRepository: dbQuerier,
	}
}

type UserService struct {
	userRepository userdb.Querier
}

func (u *UserService) CreateUser(ctx context.Context, name, password string) error {
	hashedPassword, err := u.GeneratePasswordHash(password)
	if err != nil {
		return err
	}

	createUserParams := userdb.CreateUserParams{Name: name, HashedPassword: hashedPassword}
	_, err = u.userRepository.CreateUser(ctx, createUserParams)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) LoginUser(ctx context.Context, name string, password string) (userdb.User,error) {
	user, err := u.userRepository.GetUserByName(ctx, name)
	if err != nil {
		return userdb.User{},err
	}
	err = u.CheckHashedPassword(password, user.HashedPassword)
	if errors.Is(err,bcrypt.ErrMismatchedHashAndPassword){
		return userdb.User{},customErrors.ErrMismatchedPassword
	}
	if err != nil{
		return userdb.User{},err
	}
	return user,nil

}

func (u *UserService) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err

	}
	return string(bytes), nil
}

func (u *UserService) CheckHashedPassword(password string, hashedPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err

}
