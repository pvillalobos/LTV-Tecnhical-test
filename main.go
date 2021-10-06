package main

import (
	"net/http"

	controller "github.com/bayronaz/LTV-Tecnhical-test/Controller"
	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	Utils "github.com/bayronaz/LTV-Tecnhical-test/Helpers"
	service "github.com/bayronaz/LTV-Tecnhical-test/Service"
	"github.com/gin-gonic/gin"
)

var (
	songService    service.SongService = service.New()
	songController controller.SongController
)

func main() {

	server := gin.New()

	server.Use(gin.Recovery(), gin.Logger())
	gin.SetMode(gin.ReleaseMode)

	server.GET("/releases", func(ctx *gin.Context) {
		from, until, artist, err := Utils.GetParameters(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
			return
		}
		songService = service.New()
		songController = controller.New(songService, artist)

		songController.ProcessReleasesRequest(ctx, from, until)
	})

	server.Run("localhost:8081")
}
