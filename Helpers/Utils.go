package Utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const Parse_Layout string = "2006-01-02"
const Parse_Layout_MM string = "2006-01"
const api_key string = "ec093dd5-bbe3-4d8e-bdac-314b40afb796"

var Cache = cache.New((24*time.Hour)*30, (24*time.Hour)*30)

//Read query string params from gin context
func GetParameters(ctx *gin.Context) (time.Time, time.Time, string, error) {
	var errorMsg = ""

	artist := ctx.Query("artist")

	from, err := time.Parse(Parse_Layout, ctx.Query("from"))
	if err != nil && errorMsg == "" {
		errorMsg = "Invalid or missing 'From' parameter"
	}

	until, err := time.Parse(Parse_Layout, ctx.Query("until"))
	if err != nil && errorMsg == "" {
		errorMsg = "Invalid or missing 'Until' parameter"
	}

	if from.After(until) && errorMsg == "" {
		errorMsg = "'from' is greater than 'until'"
	}
	if errorMsg != "" {
		return from, until, artist, errors.New(errorMsg)
	} else {
		return from, until, artist, nil
	}

}

//Responsible to create a Date dimension
func RangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

//Responsible for consuming API base on mode (monthly or daily) and process the response
func ConsumeSongsRepositoryAPI(date string, mode string, ctx *gin.Context) string {
	url := fmt.Sprintf("https://de-challenge.ltvco.com/v1/songs/%v?api_key=%v&released_at%v=", mode, api_key, date)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode == http.StatusBadRequest {

		var reqError entity.ErrorResponse
		var jsoninput = []byte(string(body))
		errObj := json.Unmarshal(jsoninput, &reqError)
		if errObj != nil {
			log.Fatalln("could not Unmarshal:", errObj)
		}
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
		ctx.Abort()
		return ""
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
		ctx.Abort()
		return ""
	}

	return string(body)
}
