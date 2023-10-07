package users

import "github.com/taaanechka/rest-api-go/internal/api-server/services/ports/userstorage"

type GetUserResp struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func convertBLToGetUserResp(uDB *userstorage.User) GetUserResp {
	return GetUserResp{
		ID:       uDB.ID,
		Email:    uDB.Email,
		Username: uDB.Username,
	}
}
