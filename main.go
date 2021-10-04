package main

import (
	"net/http"
	"time"

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

		getReleases(ctx, from, until, artist)
	})

	server.Run("localhost:8081")
}

// getAlbums responds with the list of all albums as JSON.
func getReleases(ctx *gin.Context, from time.Time, until time.Time, artist string) {

	for rd := Utils.RangeDate(from, until); ; {
		date := rd()

		//if func RangeDate return no dates, breaks cycle
		if date.IsZero() {
			break
		}

		//Lets look for what we have already store in caché
		songs, found := Utils.Cache.Get(date.Format(Utils.Parse_Layout))
		if !found {
			songController.AddNotFoundDates(date)
		} else {
			songController.BuildResponse(songs.([]entity.SongsRepositoryAnswer))
		}
	}

	//Check if there is missing dates to consume API
	if songController.ExistNotFoundDates() {
		songController.GetDataForNotFoundDates(ctx)
		songController.BuildResponse(nil)
	}
	res := songController.GetReleases()
	ctx.IndentedJSON(http.StatusOK, res)

}
