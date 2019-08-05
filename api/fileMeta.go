package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
	"github.com/kuuyee/matryoshka-b-multimedia/model"
)

// FileMeta 获取文件信息
// swagger:operation GET /rest/media/{type}/{ident}/meta media retrieveFileMeta
//
// 获取文件的基本信息.
//
// ---
// consumes: [application/json]
// produces: [application/json]
// security: []
// parameters:
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
func (a *API) FileMeta(c *gin.Context) {
	a.withServiceHandler(c, func(h handlers.H) {
		a.withIdent(c, func(ident string) {
			param := queryToParams(c.Request.URL.Query())
			fileOutput, size, _, err := h.RetrieveData(ident, param)
			if fileOutput != nil {
				defer fileOutput.Close()
			}
			if err != nil {
				c.AbortWithError(404, err)
				return
			}
			c.JSON(200, model.Meta{
				Ident:  ident,
				Length: size,
				Type:   h.Type(),
			})
		})
	})
}
