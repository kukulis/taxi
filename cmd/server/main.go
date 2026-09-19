package main

import (
	"darbelis.eu/taxi/internal/di"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	di.RegisterWebRoutes(router)
	//
	//// Swagger UI
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//
	//// API endpoints
	//apiRoute := router.Group("/api")
	//

	//di.RegisterApiRoutes(apiRoute)
	//
	//go di.GetHubInstance().Run()
	//di.RegisterWsRoutes(router)
	//
	_ = router.Run(":8880")
}
