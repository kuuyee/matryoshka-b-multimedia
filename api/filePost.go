package api

import (
	"errors"
	"io"

	"github.com/h2non/filetype"

	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/mime"
	"github.com/kuuyee/matryoshka-b-multimedia/model"
)

// PostFile 上传文件
// swagger:operation POST /rest/media/{type} media postFile
//
// 上传文件.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: []
// parameters:
// - name: query
//   in: query
//   description: 处理参数，另行文档
//   required: false
// - name: type
//   in: path
//   description: 文件类型
//   required: true
//   type: string
// - name: file
//   in: formData
//   description: 文件内容
//   required: true
//   type: file
// responses:
//   200:
//     description: OK
//     schema:
//         $ref: "#/definitions/Meta"
//   400:
//     description: Bad Request
//     schema:
//         $ref: "#/definitions/Error"
//   401:
//     description: Unauthorized
//     schema:
//         $ref: "#/definitions/Error"
//   403:
//     description: Forbidden
//     schema:
//         $ref: "#/definitions/Error"
//   404:
//     description: Not Found
//     schema:
//         $ref: "#/definitions/Error"
//   500:
//     description: Internal Server Error
//     schema:
//         $ref: "#/definitions/Error"
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
		c.JSON(200, model.Meta{
			Ident: ident,
			Type:  h.Type(),
		})
	})
}
