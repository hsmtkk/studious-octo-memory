package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Address string
	Message string
	UserID  int
	GroupID int
}
