package main

import (
	"time"
)

type Secret struct {
	Time   time.Time
	Secret string
}

type User struct {
	Username string
	Discord  string `form:"discord"`
	IALab    string `form:"ialab"`
	Email    string `form:"email"`
	HashedPw string
	Notes    string `form:"notes"`
	Admin    bool
	Disabled bool
}

type Event struct {
	ID           uint
	Title        string    `form:"title"`
	SignupsOpen  time.Time `form:"signupsopen"`
	SignupsClose time.Time `form:"signupsclose"`
	EventStart   time.Time `form:"eventstart"`
	EventEnd     time.Time `form:"eventend"`
	VApp         string    `form:"vapp"`
	Notes        string    `form:"notes"`
	LinkOne      string    `form:"linkone"`
	LinkTwo      string    `form:"linktwo"`
}

type SignUp struct {
	ID      uint
	User    string
	EventID uint
}

type DeployRequest struct {
	VApp     string   `json:"vapp"`
	Variants []string `json:"variants"`
}

func (u *User) IsValid() bool {
	return u.Username != "" && !u.Disabled
}
