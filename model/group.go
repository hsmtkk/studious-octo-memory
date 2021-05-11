package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	UserID  int
	Name    string
	Message string
}
