package entity

type Participant struct {
	Id     int
	UserId int

	Write   bool
	Post    bool
	Comment bool
	Delete  bool

	AddParticipant bool
}
