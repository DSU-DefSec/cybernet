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
	VApp         string    `form:"vapp"`
	VAppMandatory bool     `form:"mandatory"`
	SignupsOpen  time.Time `form:"signupsopen"`
	SignupsClose time.Time `form:"signupsclose"`
	EventStart   time.Time `form:"eventstart"`
	EventEnd     time.Time `form:"eventend"`
	Notes        string    `form:"notes"`
	LinkOne      string    `form:"linkone"`
	LinkTwo      string    `form:"linktwo"`
}

type SignUp struct {
	ID      uint
	Time    time.Time
	User    string
	EventID uint
}

type DeployRequest struct {
	Template  string   `json:"template"`
	Catalog   string   `json:"catalog"`
	Variants  []string `json:"variants"`
	MakeOwner bool     `json:"make_owner"`
}

type Attendance struct {
	Event     string     `json:"event"`
	Attendees []Attendee `json:"attendees"`
}

type Attendee struct {
	Time     time.Time `json:"time"`
	Username string    `json:"username"`
}

func (u *User) IsValid() bool {
	return u.Username != "" && !u.Disabled
}
