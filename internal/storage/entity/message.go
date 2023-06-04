package entity

import "time"

type Message struct {
	Id     int
	UserId int

	Changed bool

	Text               string
	CommentedMessageId int

	Created time.Time
}
