package middlewares

import (
	"fmt"
	"net/http"

	entity "github.com/bayronaz/LTV-Tecnhical-test/Entities"
	"github.com/gin-gonic/gin"
)

func HandleErrors(c *gin.Context) {
	c.Next() // execute all the handlers

	// at this point, all the handlers finished. Let's read the errors!
	if len(c.Errors) == 0 {
		return
	}

	err := c.Errors.Last()
	// .. create customErrorObj using err

	c.AbortWithStatusJSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
	fmt.Println(c.Writer.Header().Get("Content-Type"))
}
