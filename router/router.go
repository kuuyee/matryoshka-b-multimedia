package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-im/model"
	"github.com/kuuyee/matryoshka-b-multimedia/api"
)

func errorHandler(c *gin.Context) {
	c.Header("content-type", "application/json")

	c.Next()

	if len(c.Errors) > 0 {
		errString := strings.Join(c.Errors.Errors(), "; ")
		errCode := 500
		if existingCode := c.Writer.Status(); existingCode >= 400 {
			errCode = existingCode
		}

		c.JSON(errCode, &model.Error{
			Code:    errCode,
			Message: errString,
		})
	}
}

func New(api *api.API) *gin.Engine {
	g := gin.Default()

	mediaRoute := g.Group("/media/")
	{
		mediaRoute.Use(errorHandler)
		mediaRoute.GET(":service/:ident", api.RetrieveFile)
		mediaRoute.POST(":service", api.PostFile)
	}

	return g
}
