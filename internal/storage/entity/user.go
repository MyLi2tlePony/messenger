package entity

type User struct {
	Id int

	PublicId string

	Login    string
	Password string

	FirstName  string
	SecondName string

	Created string
}

func (u *User) Equals(user User) bool {
	return u.Id == user.Id &&
		u.PublicId == user.PublicId &&
		u.Login == user.Login &&
		u.Password == user.Password &&
		u.FirstName == user.FirstName &&
		u.SecondName == user.SecondName
}
