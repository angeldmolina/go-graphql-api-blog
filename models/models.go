package models

import (
	"time"
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	ID      uint      `json:"id" gorm:"primaryKey`
	Name    string      `json:"name" gorm:"text;not null`

	Posts []Post `json: "posts"`
}

type Post struct {
	gorm.Model
	ID        uint  `json:"id" gorm:"primaryKey`
	Title     string  `json:"title" gorm:"text;not null;<-:update`
	Content   string  `json:"content" gorm:"text;not null;<-:update`
	CreatedAt time.Time `json: "publishedAt"`
  	UpdatedAt time.Time `json: "updatedAt"`

	Author Author `json: "author"`
}
