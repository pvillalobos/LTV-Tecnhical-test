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
	return from, until, artist, errors.New(errorMsg)
}
