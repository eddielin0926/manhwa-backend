package utils

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"manhwa/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database connection instance
var DB *gorm.DB

// Connect to database
func ConnectToDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Fail to connect the database")
	}
}

func MigrateDB() {
	DB.AutoMigrate(&models.Page{}, &models.Chapter{}, &models.Book{})
}

func CreateBooks() {
	dir_path := "./books"

	dirs, err := os.ReadDir(dir_path)
	if err != nil {
		log.Println("Cannot read file path: " + dir_path)
		return
	}
	dirs = Filter(dirs, func(e os.DirEntry) bool { return e.IsDir() })

	for _, dir := range dirs {
		manhwa_path := filepath.Join(dir_path, dir.Name())
		manhwa, err := os.ReadDir(manhwa_path)
		if err != nil {
			log.Println("Cannot read file path: " + manhwa_path)
			continue
		}
		manhwa = Filter(manhwa, func(img os.DirEntry) bool { return !img.IsDir() && img.Name() != "cover.jpg" })

		var chapters []models.Chapter
		var pages []models.Page
		for _, img := range manhwa {
			list := strings.FieldsFunc(img.Name(), func(r rune) bool { return r == '_' || r == '.' })
			ch, _ := strconv.Atoi(list[0])
			pg, _ := strconv.Atoi(list[1])
			file, _ := os.ReadFile(filepath.Join(manhwa_path, img.Name()))
			pages = append(pages, models.Page{Number: uint(pg), Content: file})
			if ch > len(chapters) {
				chapters = append(chapters, models.Chapter{Number: uint(ch), Pages: pages})
			}
		}

		var book models.Book
		DB.FirstOrCreate(&book, models.Book{
			Title: dir.Name(),
		})

		for _, ch := range chapters {
			ch.BookID = book.ID
			DB.FirstOrCreate(&ch, ch)
		}
	}
}
