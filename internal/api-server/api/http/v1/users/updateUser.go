package users

import "github.com/taaanechka/rest-api-go/internal/api-server/services/ports/userstorage"

type UserUpdReq struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func convertUpdToBL(u *UserUpdReq) userstorage.User {
	return userstorage.User{
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
	}
}

type UserPatchReq struct {
	Email        string `json:"email,omitempty"`
	Username     string `json:"username,omitempty"`
	PasswordHash string `json:"password,omitempty"`
}

func convertPatchToBL(u *UserPatchReq) userstorage.User {
	return userstorage.User{
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
	}
}
