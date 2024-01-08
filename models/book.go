package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title    string `gorm:"unique"`
	Chapters []Chapter
}

type Chapter struct {
	gorm.Model
	Number uint
	Pages  []Page
	BookID uint
}

type Page struct {
	gorm.Model
	Number    uint
	Content   []byte
	ChapterID uint
}
