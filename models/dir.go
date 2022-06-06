package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Dir 收藏夹
type Dir struct {
	gorm.Model
	UserId  int64  `json:"user_id" gorm:"column:user_id"`
	DirId   int64  `json:"dir_id" gorm:"UNIQUE;column:dir_id"`
	DirName string `json:"dir_name" gorm:"column:dir_name"`
}
