package api

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
)

func (a *API) RetrieveFile(c *gin.Context) {
	a.withServiceHandler(c, func(h handlers.H) {
		a.withIdent(c, func(ident string) {
			param := queryToParams(c.Request.URL.Query())
			fileOutput, mime, err := h.RetrieveData(ident, param)
			if fileOutput != nil {
				defer fileOutput.Close()
			}
			if err != nil {
				c.AbortWithError(404, err)
				return
			}
			c.Header("content-type", mime)
			io.Copy(c.Writer, fileOutput)
		})
	})
}
