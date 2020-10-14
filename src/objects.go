package main

import (
	"database/sql"
	"time"
)

type User struct {
	Id           int
	FirstName    string
	LastName     string
	City         sql.NullString
	Birthday     sql.NullTime
	Avatar       sql.NullString
	PasswordHash string
}

type Day struct {
	Id     int
	UserId int
	Date   time.Time
}

type Activity struct {
	TypeId     int
	DayId      int
	Proportion int
}

type Emotion Activity

type ActivityType struct {
	Id         int
	UserId     int
	Name       string
	Color      string
	IsEveryday bool
}

type EmotionType ActivityType
