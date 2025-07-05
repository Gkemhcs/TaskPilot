package user

import (
	"context"

	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
)

// User represents the structure of a user in API requests.
type User struct {
	Name     string `json:"name" binding:"required"`     // Username
	Password string `json:"password" binding:"required"` // User password
	Email   string `json:"email" binding:"required"`  //User Email
}

// IUserService defines the interface for user-related business logic.
type IUserService interface {
	CreateUser(ctx context.Context, username, password,email  string) error
	LoginUser(ctx context.Context, email, password string) (*userdb.User, error)
	
}

type RefreshRequest struct {
	RefreshToken string `json:"token"`
}
