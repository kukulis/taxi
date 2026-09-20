package main

import (
	"darbelis.eu/taxi/internal/di"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//
	//// Swagger UI
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//
	//// API endpoints
	//apiRoute := router.Group("/api")
	//

	//di.RegisterApiRoutes(apiRoute)
	//

	go di.GetDriversHub().Run()
	go di.GetPassengersHub().Run()

	di.RegisterWsRoutes(router)
	di.RegisterWebRoutes(router)

	_ = router.Run(":8880")
}
