package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID  int
	PostID  int
	Message string
}
