package api

import (
	v1 "manhwa/api/v1"

	"github.com/gin-gonic/gin"
)

func SetupApi(app *gin.Engine) {
	apiRoute := app.Group("api")
	v1.AddApi(apiRoute)
}
