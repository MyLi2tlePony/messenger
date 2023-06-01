package entity

type Participant struct {
	UserId int

	Write   bool
	Delete  bool
	Read    bool
	Comment bool
}
