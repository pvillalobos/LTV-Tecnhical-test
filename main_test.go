package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "github.com/bayronaz/LTV-Tecnhical-test/Controller"
	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	Utils "github.com/bayronaz/LTV-Tecnhical-test/Helpers"
	service "github.com/bayronaz/LTV-Tecnhical-test/Service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Router() *gin.Engine {
	server := gin.New()
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
	return server
}

func Test_ProcessReleasesRequest(t *testing.T) {
	fmt.Println("")
	fmt.Println("-----> Test_ProcessReleasesRequest")
	request, _ := http.NewRequest("GET", "/releases?from=2021-01-01&until=2021-01-05", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	//fmt.Println(response.Body)
	assert.Equal(t, 200, response.Code, "OK Response is expected")
}

func Test_MissingFromParameter(t *testing.T) {
	fmt.Println("")
	fmt.Println("-----> Test_MissingFromParameter")
	request, _ := http.NewRequest("GET", "/releases?until=2021-01-05", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	fmt.Println(response.Body)
	assert.Equal(t, 400, response.Code, "400 StatusCode response is expected")
}

func Test_MissingUntilParameter(t *testing.T) {
	fmt.Println("")
	fmt.Println("-----> Test_MissingUntilParameter")
	request, _ := http.NewRequest("GET", "/releases?from=2021-01-01", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	fmt.Println(response.Body)
	assert.Equal(t, 400, response.Code, "400 StatusCode response is expected")
}

func Test_ProcessReleasesRequestWithArtist(t *testing.T) {
	fmt.Println("")
	fmt.Println("-----> Test_ProcessReleasesRequestWithArtist")
	request, _ := http.NewRequest("GET", "/releases?artist=Camilo&from=2021-03-01&until=2021-03-05", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	res := []entity.OutputResponse{
		{ReleasedAt: "2021-03-01",
			Songs: []entity.Songs{
				{Artist: "Camilo", Name: "Machu Picchu"},
			},
		},
		{ReleasedAt: "2021-03-04",
			Songs: []entity.Songs{
				{Artist: "Camilo", Name: "Manos de Tijera"},
				{Artist: "Camilo", Name: "Ropa Cara"},
			},
		},
	}

	body, _ := ioutil.ReadAll(response.Body)
	var bodyRes []entity.OutputResponse
	var jsoninput = []byte(string(body))
	json.Unmarshal(jsoninput, &bodyRes)

	assert.Equal(t, bodyRes, res, "OK Response is expected")
}

/*func Test_DateOutOfRange(t *testing.T) {
	request, _ := http.NewRequest("GET", "/releases?from=2020-01-01&until=2020-01-05", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	fmt.Println(response.Body)
	assert.Equal(t, 400, response.Code, "OK Response is expected")
}*/
