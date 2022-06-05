package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Collection struct {
	gorm.Model
	CollectionId int64  `json:"collection_id" gorm:"UNIQUE;column:collection_id"`
	DirId        int64  `json:"dir_id" gorm:"column:dir_id"`
	Word         string `json:"word" gorm:"word"`
	UrlName      string `json:"url_name" gorm:"url_id"`
}
