package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
)

func (a *API) withServiceHandler(c *gin.Context, next func(h handlers.H)) {
	h, ok := a.handlers[c.Param("service")]
	if !ok {
		c.AbortWithError(404, errors.New("service not found"))
		return
	}
	next(h)
}
