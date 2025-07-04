package types

type IUserHandler interface{
	CreateUser()
	LoginUser()
	CheckHashedPassword()
}





