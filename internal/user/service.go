package user

import (
	"context"
	"errors"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// NewUserService creates a new UserService with the given user repository (database querier).
func NewUserService(dbQuerier userdb.Querier) *UserService {
	return &UserService{
		userRepository: dbQuerier,
	}
}

// UserService implements user-related business logic and interacts with the database.
type UserService struct {
	userRepository userdb.Querier // Database access layer for user operations
}

// CreateUser hashes the password and creates a new user in the database.
func (u *UserService) CreateUser(ctx context.Context, name, password string) error {
	hashedPassword, err := u.GeneratePasswordHash(password)
	if err != nil {
		return err
	}

	createUserParams := userdb.CreateUserParams{Name: name, HashedPassword: hashedPassword}
	_, err = u.userRepository.CreateUser(ctx, createUserParams)
	if err, ok := err.(*pq.Error); ok {
	if err.Code == "23505" {
		// Check constraint name if needed
		if err.Constraint == "users_pkey" {
			return customErrors.ErrUserAlreadyExists
		}
	}
}
	if err != nil {
		return err
	}
	return nil
}

// LoginUser authenticates a user by name and password.
// Returns the user if authentication is successful.
func (u *UserService) LoginUser(ctx context.Context, name string, password string) (userdb.User, error) {
	user, err := u.userRepository.GetUserByName(ctx, name)
	if err != nil {
		return userdb.User{}, err
	}
	// Compare the provided password with the stored hashed password
	err = u.CheckHashedPassword(password, user.HashedPassword)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return userdb.User{}, customErrors.ErrMismatchedPassword
	}
	if err != nil {
		return userdb.User{}, err
	}
	return user, nil
}

// GeneratePasswordHash hashes a plain-text password using bcrypt.
func (u *UserService) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckHashedPassword compares a plain-text password with a hashed password.
func (u *UserService) CheckHashedPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
