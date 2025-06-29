package domain

import "time"

type Account struct {
	ID             uint64
	Name           string
	Currency       string
	AccountType    string
	InitialBalance float64
	PrimaryOwnerID uint64
	SecondOwnerID  *uint64 // Optional - pointer to allow nil
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
