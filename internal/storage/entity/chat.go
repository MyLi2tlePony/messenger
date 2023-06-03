package entity

import "time"

type Chat struct {
	Id int

	Description string
	Name        string

	Open    bool
	Created time.Time
}
