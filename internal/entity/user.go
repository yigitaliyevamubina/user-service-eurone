package entity

import "time"

type User struct {
	Id           string
	Username     string
	Email        string
	Password     string
	FirstName    string
	LastName     string
	Bio          string
	Website      string
	IsActive     bool
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type GetListFilter struct {
	Page    int64  `json:"page"`
	Limit   int64  `json:"limit"`
	OrderBy string `json:"order_by"`
}