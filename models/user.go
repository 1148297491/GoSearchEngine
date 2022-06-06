package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	UserID    int64  `json:"user_id"`
	Userphone string `json:"userphone"`
	Password  string `json:"password"`
}

type TUser struct {
	UserID    string
	Userphone string
}
