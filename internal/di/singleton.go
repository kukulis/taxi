package di

import "darbelis.eu/taxi/internal/web"

var webControllerInstance *web.WebController

func GetWebController() *web.WebController {
	if webControllerInstance == nil {
		webControllerInstance = web.NewWebController()
	}

	return webControllerInstance
}
