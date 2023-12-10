package service

import (
	"errors"
	"log"
	"rest_app/database"

	"github.com/gofrs/uuid"
)

type UserInfo struct {
	Id    uuid.UUID
	Email string
}

type User interface {
	FindById(userId uuid.UUID) (*UserInfo, error)
	FindByEmail(email string) (*UserInfo, error)
}
type UserService struct {
	Db database.Database
}

func (u UserService) FindById(userId uuid.UUID) (*UserInfo, error) {
	var userInfo = UserInfo{}
	var err = u.Db.QueryRow(
		"SELECT * FROM users WHERE id = $1",
		userId,
	)(&userInfo.Id, &userInfo.Email)

	if err != nil {
		return nil, errors.New("user is not found")
	}

	return &userInfo, nil
}

func (u UserService) FindByEmail(email string) (*UserInfo, error) {
	var userInfo = UserInfo{}
	var err = u.Db.QueryRow(
		"SELECT * FROM users WHERE email = $1",
		email,
	)(&userInfo.Id, &userInfo.Email)

	if err != nil {
		log.Println("an error during getting user's information: ", err)
		return nil, errors.New("user is not found")
	}

	return &userInfo, nil
}
