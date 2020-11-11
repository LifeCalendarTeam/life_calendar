package main

import (
	"database/sql"
	"time"
)

// Logical objects:

// User object describes a user of LifeCalendar (identifier, profile info, password hash)
type User struct {
	Id           int            `db:"id"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"second_name"`
	City         sql.NullString `db:"city"`
	Birthday     sql.NullTime   `db:"birthday"`
	Avatar       sql.NullString `db:"avatar"`
	PasswordHash string         `db:"password_hash"`
}

// Day object describes a day (identifier, user identifier, date)
type Day struct {
	Id     int       `db:"id"`
	UserId int       `db:"user_id"`
	Date   time.Time `db:"date"`
}

// ActivityType describing describes an activity type (identifier, user identifier, name/label, color, is it everyday)
type ActivityType struct {
	Id         int    `db:"id"`
	UserId     int    `db:"user_id"`
	Name       string `db:"name"`
	Color      string `db:"color"`
	IsEveryday bool   `db:"is_everyday"`
}

// EmotionType object describes an emotion type (semantically similar to `ActivityType`)
type EmotionType ActivityType

// Activity object describes an activity (type identifier, day identifier, proportion value)
type Activity struct {
	TypeId     int `db:"type_id"`
	DayId      int `db:"day_id"`
	Proportion int `db:"proportion"`
}

// Emotion object describes an emotion (semantically similar to `Activity`)
type Emotion Activity

// Forms:

type loginForm struct {
	UserId   int    `schema:"user_id,required"`
	Password string `schema:"password,required"`
}

// Internal objects:

type proportionAndColor struct {
	Proportion float64 `db:"proportion"`
	Color      string  `db:"color"`
}

type briefDay struct {
	DayId        int       `db:"id" json:"id"`
	Date         time.Time `db:"date" json:"date"`
	AverageColor [3]int    `json:"average_color"`
}
