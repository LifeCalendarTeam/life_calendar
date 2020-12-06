package main

import (
	"database/sql"
	"time"
)

type User struct {
	Id           int            `db:"id"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"second_name"`
	City         sql.NullString `db:"city"`
	Birthday     sql.NullTime   `db:"birthday"`
	Avatar       sql.NullString `db:"avatar"`
	PasswordHash string         `db:"password_hash"`
}

type Day struct {
	Id     int       `db:"id"`
	UserId int       `db:"user_id"`
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
