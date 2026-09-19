package di

import "github.com/gin-gonic/gin"

func RegisterWebRoutes(router *gin.Engine) {

	// Templates
	router.LoadHTMLFiles(
		"./pages/main.gohtml",
		"./pages/driver.gohtml",
		"./pages/passenger.gohtml",
	)

	router.GET("/", func(c *gin.Context) { GetWebController().Index(c) })
	router.GET("/driver", func(c *gin.Context) { GetWebController().Driver(c) })
	router.GET("/passenger", func(c *gin.Context) { GetWebController().Passenger(c) })

	router.Static("/assets", "./assets")

	router.StaticFile("/favicon.ico", "./assets/img/favicon.ico")
}
