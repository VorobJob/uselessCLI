package models

import "time"

type Records struct {
	ID uint `gorm:"primary key;autoIncrement"`
	User
}

type User struct {
	Name *string   `gorm:"not null"`
	DOB  time.Time `gorm:"not null"`
	Sex  bool      `gorm:"not null"`
}
