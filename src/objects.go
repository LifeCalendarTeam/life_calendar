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

// User object describes a user of LifeCalendar (identifier, profile info, password hash)
type User struct {
	ID        int            `db:"id"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"second_name"`
	City      sql.NullString `db:"city"`
	Birthday  sql.NullTime   `db:"birthday"`
	Avatar    sql.NullString `db:"avatar"`
}

// Day object describes a day (identifier, user identifier, date)
type Day struct {
	ID     int       `db:"id"`
	UserID int       `db:"user_id"`
	Date   time.Time `db:"date"`
}

// ActivityOrEmotion object describes an activity/emotion (type identifier, day identifier, proportion value)
type ActivityOrEmotion struct {
	TypeID     int `db:"type_id" json:"type_id,omitempty"`
	DayID      int `db:"day_id" json:"day_id,omitempty"`
	Proportion int `db:"proportion" json:"proportion,omitempty"`
}

// ActivityOrEmotionType describes an activity/emotion type (identifier, user identifier, name/label, color, is it
// everyday)
type ActivityOrEmotionType struct {
	Id         int                      `db:"id"`
	UserId     int                      `db:"user_id"`
	Name       string                   `db:"name"`
	Color      string                   `db:"color"`
	IsEveryday bool                     `db:"is_everyday"`
	EntityType WhetherActivityOrEmotion `db:"activity_or_emotion"`
}

// Internal objects

type activityOrEmotionWithType struct {
	ActivityOrEmotion
	EntityType WhetherActivityOrEmotion `db:"activity_or_emotion"`
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
