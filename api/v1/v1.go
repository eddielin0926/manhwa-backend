package v1

import (
	"manhwa/api/v1/books"

	"github.com/gin-gonic/gin"
)

func AddApi(rg *gin.RouterGroup) {
	v1 := rg.Group("v1")
	books.AddBookRoute(v1)
}
