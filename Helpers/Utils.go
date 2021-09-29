package Utils

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

const Parse_Layout string = "2006-01-02"

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
func rangeDate(start, end time.Time) func() time.Time {
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
