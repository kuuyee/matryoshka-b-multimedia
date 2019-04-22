package api

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
)

// RetrieveFile 下载文件
// swagger:operation GET /rest/media/{type}/{ident} media retrieveFile
//
// 下载文件.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: []
// parameters:
// - name: query
//   in: query
//   description: 处理参数，另行文档
// - name: type
//   in: path
//   description: 文件类型
//   required: true
//   type: string
// - name: ident
//   in: path
//   description: 文件识别符
//   required: true
//   type: string
// responses:
//   200:
//     description: OK
//     schema:
//       type: file
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
