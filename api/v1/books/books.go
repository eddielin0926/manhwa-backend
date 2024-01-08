package books

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"manhwa/models"
	"manhwa/utils"

	"github.com/gin-gonic/gin"
)

func AddBookRoute(rg *gin.RouterGroup) {
	books := rg.Group("books")
	books.GET("", getBooks)
	books.POST("", creatBook)
	books.GET("/:title/cover", getBookCover)
	books.GET("/:title/contents/:chapter/:page", getBookContent)
}

// GetBooks 	 godoc
//
//	@Summary		Get a list of books' title
//	@Description	Return a list of books' title
//	@Tags			Books
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	string
//	@Router			/books [get]
func getBooks(ctx *gin.Context) {
	var books []string
	utils.DB.Table("books").Select("title").Find(&books)
	ctx.JSON(http.StatusOK, gin.H{"books": books})
}

type BookForm struct {
	Title    string                  `form:"title" binding:"required"`
	Cover    *multipart.FileHeader   `form:"cover" binding:"required"`
	Contents []*multipart.FileHeader `form:"contents" binding:"required"`
}

// CreateBook godoc
//
//	@Summary		Create book
//	@Description	Create book
//	@ID				books.create-book
//	@Tags			Books
//	@Accept			multipart/form-data
//	@Param			title		formData	string	true	"Book Title"
//	@Param			cover		formData	file	true	"Book Cover"
//	@Param			contents	formData	[]file	true	"Book Contents"
//	@Success		200
//	@Router			/books [post]
func creatBook(ctx *gin.Context) {
	var form BookForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	var chapters []models.Chapter
	var pages []models.Page
	for _, contentFile := range form.Contents {
		openedFile, _ := contentFile.Open()
		file, _ := io.ReadAll(openedFile)
		list := strings.FieldsFunc(contentFile.Filename, func(r rune) bool { return r == '_' || r == '.' })
		ch, _ := strconv.Atoi(list[0])
		pg, _ := strconv.Atoi(list[1])
		pages = append(pages, models.Page{Number: uint(pg), Content: file})
		if ch > len(chapters) {
			chapter := models.Chapter{
				Number: uint(ch),
				Pages:  pages,
			}
			chapters = append(chapters, chapter)
			pages = []models.Page{}
		}
	}
	book := models.Book{
		Title:    form.Title,
		Chapters: chapters,
	}
	utils.DB.Create(&book)
	ctx.String(http.StatusOK, "Uploaded successfully")
}

// GetBookCover 	 godoc
//
//	@Summary		Get book cover
//	@Description	get book's cover image
//	@Tags			Books
//	@Accept			json
//	@Produce		jpeg
//	@Param			title	path	string	true	"Book Title"
//	@Success		200		{file}	jpeg
//	@Router			/books/{title}/cover [get]
func getBookCover(ctx *gin.Context) {
	title := []byte(ctx.Param("title"))

	filePath := fmt.Sprintf("./books/%s/cover.jpg", title)

	ctx.File(filePath)
}

// GetBookCover 	 godoc
//
//	@Summary		Get book content
//	@Description	get book's content
//	@Tags			Books
//	@Accept			json
//	@Produce		jpeg
//	@Param			title	path	string	true	"Book Title"
//	@Param			chapter	path	int		true	"Book Chapter"
//	@Param			page	path	int		true	"Book Page"
//	@Success		200		{file}	jpeg
//	@Router			/books/{title}/contents/{chapter}/{page} [get]
func getBookContent(ctx *gin.Context) {
	title := []byte(ctx.Param("title"))
	chapter, _ := strconv.Atoi(ctx.Param("chapter"))
	page, _ := strconv.Atoi(ctx.Param("page"))

	filePath := fmt.Sprintf("./books/%s/%02d_%02d.jpg", title, chapter, page)

	ctx.File(filePath)
}
