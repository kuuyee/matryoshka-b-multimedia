package router

import (
	"log"
	"bytes"
	"errors"
	"strings"
	"net/http"

	"github.com/gobuffalo/packr/v2"
	
	"github.com/gin-contrib/location"

	"github.com/gin-gonic/gin"

	"github.com/jmattheis/go-packr-swagger-ui"
	
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
	g.Use(location.Default())

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

	// 注册swagger docs
	if gin.Mode() == gin.DebugMode {
		swaggerUIBox := swaggerui.GetBox()
		docs := packr.New("swagger-ui", "../docs/")
		g.GET("/docs/*any", gin.WrapH(http.StripPrefix("/docs/", http.FileServer(swaggerUIBox))))
		g.GET("/swagger", func(c *gin.Context) {
			swaggerDef, err := docs.Find("spec.json")
			if err != nil {
				c.AbortWithError(404, errors.New("spec file not found"))
				return
			}
			loc := location.Get(c)
			log.Println(loc)
			swaggerDef = bytes.Replace(swaggerDef, []byte("{{IM_REST_HOST}}"), []byte(loc.Host), 1)

			c.Header("content-type", "application/json")
			c.Writer.Write(swaggerDef)
		})
	}

	return g
}
