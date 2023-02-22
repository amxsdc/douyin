package model

import "github.com/jinzhu/gorm"

type Message struct {
	gorm.Model
	ToUserId   uint   `json:"to_user_id"`
	FromUserId uint   `json:"from_user_id"`
	Content    string `json:"content"`
}
