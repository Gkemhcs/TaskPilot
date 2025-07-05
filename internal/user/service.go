package user

import (
	"context"
	"database/sql"
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
func (u *UserService) CreateUser(ctx context.Context, name,email, password string) error {
	hashedPassword, err := u.GeneratePasswordHash(password)
	if err != nil {
		return err
	}

	createUserParams := userdb.CreateUserParams{Name: name, HashedPassword: hashedPassword,Email:email}
	_, err = u.userRepository.CreateUser(ctx, createUserParams)
	
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" {
			if pqErr.Constraint == "users_pkey" {
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
func (u *UserService) LoginUser(ctx context.Context, email string, password string) (*userdb.User, error) {
	user, err := u.userRepository.GetUserByEmail(ctx,email)
	
	if errors.Is(err,sql.ErrNoRows){
		return nil,customErrors.USER_NOT_FOUND
	}
	if err != nil {
		return nil, err
	}
	
	// Compare the provided password with the stored hashed password
	err = u.CheckHashedPassword(password, user.HashedPassword)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, customErrors.ErrMismatchedPassword
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
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
