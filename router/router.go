package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/api"
)

func New(api *api.API) *gin.Engine {
	g := gin.Default()

	mediaRoute := g.Group("/media/")
	{
		mediaRoute.GET(":service/:ident", api.RetrieveFile)
		mediaRoute.POST(":service", api.PostFile)
	}

	return g
}
