package mongodb

import "github.com/taaanechka/rest-api-go/internal/api-server/services/ports/userstorage"

type User struct {
	ID           string `bson:"_id,omitempty"`
	Email        string `bson:"email"`
	Username     string `bson:"username"`
	PasswordHash string `bson:"password,omitempty"`
}

func convertDBToBL(uDB *User) userstorage.User {
	return userstorage.User{
		ID:           uDB.ID,
		Email:        uDB.Email,
		Username:     uDB.Username,
		PasswordHash: uDB.PasswordHash,
	}
}

func convertBLToDB(uBL *userstorage.User) User {
	return User{
		ID:           uBL.ID,
		Email:        uBL.Email,
		Username:     uBL.Username,
		PasswordHash: uBL.PasswordHash,
	}
}
