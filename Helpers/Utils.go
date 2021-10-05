package Utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

//Responsible to create a Date dimension array func to check dates stores on cache or to request them to song repository API
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

//Responsible for consuming API base on a mode (monthly or daily) and return the response body
func ConsumeSongsRepositoryAPI(date string, mode string, ctx *gin.Context, errorChannel chan string) string {
	url := fmt.Sprintf("https://de-challenge.ltvco.com/v1/songs/%v?api_key=%v&released_at=%v", mode, api_key, date)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)
	if err != nil {
		//ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		errorChannel <- err.Error()
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode == http.StatusBadRequest {

		var reqError entity.ErrorResponse
		var jsoninput = []byte(string(body))
		errObj := json.Unmarshal(jsoninput, &reqError)
		if errObj != nil {
			errorChannel <- "could not Unmarshal:"
			//log.Fatalln("could not Unmarshal:", errObj)
		}
		//ctx.AbortWithStatusJSON(http.StatusBadRequest, reqError)
		errorChannel <- reqError.Error
		return ""
	}

	if err != nil {
		//ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
		errorChannel <- err.Error()
	}

	return string(body)
}
