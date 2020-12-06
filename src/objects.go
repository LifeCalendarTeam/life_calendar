package main

import (
	"database/sql"
	"time"
)

// Helpers

// WhetherActivityOrEmotion is technically a synonym for `string`, but the objects of this type should only take values
// from the set of the following constants: `EntityTypeActivity`, `EntityTypeEmotion`
type WhetherActivityOrEmotion string

const (
	// EntityTypeActivity is a constant representing the "activity" type of an object which can either be an activity or
	// en emotion
	EntityTypeActivity WhetherActivityOrEmotion = "activity"
	// EntityTypeEmotion is a constant representing the "emotion" type of an object which can either be an activity or
	// an emotion
	EntityTypeEmotion WhetherActivityOrEmotion = "emotion"
)

// Logical objects

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

type ActivityOrEmotion struct {
	TypeId     int `db:"type_id"`
	DayId      int `db:"day_id"`
	Proportion int `db:"proportion"`
}

type ActivityOrEmotionType struct {
	Id         int                      `db:"id"`
	UserId     int                      `db:"user_id"`
	Name       string                   `db:"name"`
	Color      string                   `db:"color"`
	IsEveryday bool                     `db:"is_everyday"`
	EntityType WhetherActivityOrEmotion `db:"activity_or_emotion"`
}
