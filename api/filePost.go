package api

import (
	"errors"
	"io"

	"github.com/h2non/filetype"

	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/mime"
)

func (a *API) PostFile(c *gin.Context) {
	a.withServiceHandler(c, func(h handlers.H) {
		formFile, err := c.FormFile("file")
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		if formFile.Size > h.SizeLimit() {
			c.AbortWithError(413, errors.New("file size exceeded limit"))
			return
		}

		formFileReader, err := formFile.Open()
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		defer formFileReader.Close()

		mime := mime.FileNameToMIME(formFile.Filename)
		if mime == "" {
			typ, err := filetype.MatchReader(formFileReader)
			if err != nil {
				c.AbortWithError(400, errors.New("unable to determine file type"))
				return
			}
			mime = typ.MIME.Value
			_, err = formFileReader.Seek(0, io.SeekStart)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
		}

		param := queryToParams(c.Request.URL.Query())
		ident, err := h.WriteData(formFileReader, mime, param)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, struct {
			Ident string `json:"ident"`
		}{ident})
	})
}
