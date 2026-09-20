package main

import (
	"darbelis.eu/taxi/internal/di"
	"github.com/gin-gonic/gin"
)

func main() {

	dispatcher := di.GetDispatcher()
	di.InitializeListeners(dispatcher)

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
	go di.GetDriversWsHandler().Handle(di.GetDriversHub().GetIncomingMessagesChannel())

	go di.GetPassengersHub().Run()
	go di.GetPassengersWsHandler().Handle(di.GetPassengersHub().GetIncomingMessagesChannel())

	di.RegisterWsRoutes(router)
	di.RegisterWebRoutes(router)

	_ = router.Run(":8880")
}
