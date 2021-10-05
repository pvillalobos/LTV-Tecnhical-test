package main

import (
	"net/http"

	controller "github.com/bayronaz/LTV-Tecnhical-test/Controller"
	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	Utils "github.com/bayronaz/LTV-Tecnhical-test/Helpers"
	service "github.com/bayronaz/LTV-Tecnhical-test/Service"
	"github.com/gin-gonic/gin"
)

// albums slice to seed record album data.
var albums = []entity.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var (
	songService    service.SongService = service.New()
	songController controller.SongController
)

func main() {

	server := gin.New()

	server.Use(gin.Recovery(), gin.Logger())

	server.GET("/releases", func(ctx *gin.Context) {
		from, until, artist, err := Utils.GetParameters(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
			return
		}
		songService = service.New()
		songController = controller.New(songService, artist)

		songController.ProcessRequest(ctx, from, until, artist)
	})

	server.Run("localhost:8081")
}
