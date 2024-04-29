package main

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required"`
	IsOnline bool   `json:"is_online,omitempty"`
}
