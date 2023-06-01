package entity

import "time"

type Chat struct {
	Id   int
	Type int

	Name    string
	Created time.Time
}
