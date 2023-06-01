package entity

import "time"

type Message struct {
	Id            int
	CreatedUserId int
	CreatedChatId int
	SendUserId    int

	Personal bool
	Changed  bool
	Read     bool
	Text     string
	Created  time.Time

	CommentMessageId int
}
