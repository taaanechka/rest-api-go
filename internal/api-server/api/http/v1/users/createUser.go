package users

import "github.com/taaanechka/rest-api-go/internal/api-server/services/ports/userstorage"

type CreateUserReq struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func convertCreateUserReqToBL(u *CreateUserReq) userstorage.User {
	return userstorage.User{
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
	}
}
