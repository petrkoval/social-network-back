package auth

type RegisterUserDTO = CreateUserDTO

type LoginUserDTO = CreateUserDTO

type CreateUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
