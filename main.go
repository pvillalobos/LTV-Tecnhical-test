package main

import (
	"fmt"
	"net/http"
	"time"

	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	Utils "github.com/bayronaz/LTV-Tecnhical-test/Helpers"
	"github.com/gin-gonic/gin"
)

// albums slice to seed record album data.
var albums = []entity.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {

	server := gin.New()

	server.Use(gin.Recovery(), gin.Logger())

	server.GET("/releases", func(ctx *gin.Context) {
		from, until, artist, err := Utils.GetParameters(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
			ctx.Abort()
			return
		}
		getReleases(ctx, from, until, artist)
	})

	server.Run("localhost:8081")
}

// getAlbums responds with the list of all albums as JSON.
func getReleases(c *gin.Context, from time.Time, until time.Time, artist string) {

	for rd := Utils.RangeDate(from, until); ; {
		date := rd()

		//if func RangeDate return no dates, breaks cycle
		if date.IsZero() {
			break
		}

		//Missing dates in caché that will be requested to API
		var DatesNotFound = []time.Time{}

		songs, found := Utils.C.Get(date.Format(Utils.Parse_Layout))
		if !found {
			DatesNotFound = append(DatesNotFound, date)
			fmt.Println(DatesNotFound)
		} else {
			fmt.Println(songs.([]entity.SongsRepositoryAnswer))
		}
	}

	c.IndentedJSON(http.StatusOK, albums)
}
