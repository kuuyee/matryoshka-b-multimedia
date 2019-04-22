package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/api"
	"github.com/kuuyee/matryoshka-b-multimedia/model"
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

		c.JSON(errCode, model.Error{errCode, errString})
	}
}

// New creates a new gin engine
func New(api *api.API) *gin.Engine {
	g := gin.Default()

	restRoute := g.Group("/rest/")
	{
		restRoute.Use(errorHandler)
		mediaRoute := restRoute.Group("/media/")
		{
			mediaRoute.GET(":service/:ident", api.RetrieveFile)
			mediaRoute.GET(":service/:ident/meta", api.FileMeta)
			mediaRoute.POST(":service", api.PostFile)
		}
	}

	return g
}
