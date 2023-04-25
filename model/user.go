package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"column:name;type:varchar(20);not null"`
	Telephone string `gorm:"column:telephone;type:varchar(11);not null;unique"`
	Password  string `gorm:"column:password;size:255;not null"`
}
