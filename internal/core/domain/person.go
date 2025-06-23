package domain

import "time"

type Person struct {
	ID        uint64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
