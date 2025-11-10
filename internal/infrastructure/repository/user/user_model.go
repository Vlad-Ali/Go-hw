package userrepository

import (
	userdomain "transfersservice/internal/domain/user"
	"transfersservice/internal/domain/user/object"
)

type UserModel struct {
	Username string
	Password string
	Email    string
	ID       string
}

func NewUserModel(userName string, password string, email string, id string) *UserModel {
	return &UserModel{userName, password, email, id}
}

func (u *UserModel) ToDomain() (*userdomain.User, error) {
	userID, err := object.NewUserID(u.ID)
	if err != nil {
		return nil, err
	}
	user := userdomain.NewUser(u.Username, u.Password, u.Email)
	err = user.SetID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
