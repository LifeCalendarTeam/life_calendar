package main

import (
	"database/sql"
	"time"
)

// Logical objects:

// User object describes a user of LifeCalendar (identifier, profile info, password hash)
type User struct {
	ID           int            `db:"id"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"second_name"`
	City         sql.NullString `db:"city"`
	Birthday     sql.NullTime   `db:"birthday"`
	Avatar       sql.NullString `db:"avatar"`
	PasswordHash string         `db:"password_hash"`
}

// Day object describes a day (identifier, user identifier, date)
type Day struct {
	ID     int       `db:"id"`
	UserID int       `db:"user_id"`
	Date   time.Time `db:"date"`
}

// ActivityOrEmotion object describes an activity/emotion (type identifier, day identifier, proportion value)
type ActivityOrEmotion struct {
	TypeId     int `db:"type_id"`
	DayId      int `db:"day_id"`
	Proportion int `db:"proportion"`
	// TODO: probably needs a field telling whether it's an activity or an emotion
}

// ActivityOrEmotionType describes an activity/emotion type (identifier, user identifier, name/label, color, is it
// everyday)
type ActivityOrEmotionType struct {
	Id         int    `db:"id"`
	UserId     int    `db:"user_id"`
	Name       string `db:"name"`
	Color      string `db:"color"`
	IsEveryday bool   `db:"is_everyday"`
	// TODO: probably needs a field telling whether it's an activity or an emotion
}

// Forms:

type loginForm struct {
	UserID   int    `schema:"user_id,required"`
	Password string `schema:"password,required"`
}

// Internal objects:

type proportionAndColor struct {
	Proportion float64 `db:"proportion"`
	Color      string  `db:"color"`
}

type briefDay struct {
	DayID        int       `db:"id" json:"id"`
	Date         time.Time `db:"date" json:"date"`
	AverageColor [3]int    `json:"average_color"`
}
