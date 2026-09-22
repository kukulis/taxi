package web

import (
	"log"

	"darbelis.eu/taxi/pkg/util"
	"github.com/gin-gonic/gin"
)

func RollCookie(ctx *gin.Context) (string, error) {
	cookie, err := ctx.Cookie(ClientCookieName)
	var clientId string

	if err != nil {
		log.Println("Failed to get client id from cookie", err)
		clientId, err = util.RandomHash(8)
		if err != nil {
			log.Println("Failed to create a client id", err)
			return "", err
		}

		ctx.SetCookie(ClientCookieName,
			clientId,
			10000000, // 115 days
			"/",
			"",
			false,
			false)
	} else {
		clientId = cookie
	}

	return clientId, nil
}
