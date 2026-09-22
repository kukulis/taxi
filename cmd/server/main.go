package main

import (
	"fmt"
	"os"

	"darbelis.eu/taxi/internal/di"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using default values")
	}

	dispatcher := di.GetDispatcher()
	mainState := di.GetMainState()
	di.InitializeListenersFromMainState(mainState, dispatcher)

	go mainState.HandleDedicatedEvents()

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

	go di.GetTaxiDispatcher().TickForDriversUpdates()
	go di.GetTaxiDispatcher().TickForPassengersUpdates()

	di.RegisterWsRoutes(router)
	di.RegisterWebRoutes(router)

	useTslString := os.Getenv("USE_TSL")

	useTls := useTslString == "true"

	if useTls {
		_ = router.RunTLS(":5443", "./tls/server.crt", "./tls/server.key")
	} else {
		_ = router.Run(":8880")
	}
}
