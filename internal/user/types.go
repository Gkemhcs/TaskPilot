package user

import (
	"context"

	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
)


type User struct {
	
    Name     string `json:"name" binding:"required"`
    Password string `json:"password" binding:"required"`

}

type IUserService interface {
	CreateUser(ctx context.Context,username,password string) (error)
	LoginUser(ctx context.Context,username,password string ) (userdb.User ,error)
}

