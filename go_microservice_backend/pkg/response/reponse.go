package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type ErrorResponseData struct {
	Code   int         `json:"code"`
	Err    string      `json:"error"`
	Detail interface{} `json:"detail"`
}

// success response

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: MSG[code],
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int) {
	c.JSON(code, ResponseData{
		Code:    code,
		Message: MSG[code],
		Data:    nil,
	})
}
