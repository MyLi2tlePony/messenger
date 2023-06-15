package entity

import "time"

type User struct {
	Id int

	PublicId string

	Login    string
	Password string

	FirstName  string
	SecondName string

	Created time.Time
}

func NewUser(publicId, login, password, firstName, secondName string, created time.Time) User {
	return User{
		PublicId:   publicId,
		Login:      login,
		Password:   password,
		FirstName:  firstName,
		SecondName: secondName,
		Created:    created,
	}
}
