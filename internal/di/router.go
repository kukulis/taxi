package di

import "github.com/gin-gonic/gin"

func RegisterWebRoutes(router *gin.Engine) {

	// Templates
	router.LoadHTMLFiles(
		"./pages/main.gohtml",
	)

	router.GET("/", func(c *gin.Context) { GetWebController().Index(c) })
}
