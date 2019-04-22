package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
	"github.com/kuuyee/matryoshka-b-multimedia/model"
)

func (a *API) FileMeta(c *gin.Context) {
	a.withServiceHandler(c, func(h handlers.H) {
		a.withIdent(c, func(ident string) {
			param := queryToParams(c.Request.URL.Query())
			fileOutput, _, err := h.RetrieveData(ident, param)
			if fileOutput != nil {
				defer fileOutput.Close()
			}
			if err != nil {
				c.AbortWithError(404, err)
				return
			}
			c.JSON(200, model.Meta{
				Ident: ident,
			})
		})
	})
}
