package main

import (
	"net/http"

	"manhwa/api"
	"manhwa/docs"
	"manhwa/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	utils.LoadEnv()
	utils.ConnectToDB()
	utils.MigrateDB()
	// utils.CreateBooks()
}

// @title			Go Backend
// @version		0.1
// @description	This is a sample server implement with Go.
// @schemes		http
// @contact.name	Eddie Lin
// @license.name	MIT
// @license.url	https://opensource.org/licenses/MIT
func main() {
	app := gin.Default()
	app.Use(cors.Default())

	// API Docs
	docs.SwaggerInfo.BasePath = "/api/v1"
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusPermanentRedirect, "/docs/index.html") })

	api.SetupApi(app)
	app.Run()
}
