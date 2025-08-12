package templates

// response_template.go - Template for generating response files (success, error, etc.)

// SuccessResponseTemplate - Template for generating SuccessResponse
const SuccessResponseTemplate = `package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse struct represents the structure of a success response
type SuccessResponse struct {
	Code    int         ` + "`" + `json:"code"` + "`" + `
	Message string      ` + "`" + `json:"message"` + "`" + `
	Data    interface{} ` + "`" + `json:"data,omitempty"` + "`" + `
}

// Success sends a standardized success response
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// NoContentResponse struct represents the structure of a no content response (204)
type NoContentResponse struct {
	Code    int    ` + "`" + `json:"code"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
}

// SendNoContent sends a standardized no content response
func SendNoContent(c *gin.Context, message string) {
	c.JSON(http.StatusNoContent, NoContentResponse{
		Code:    http.StatusNoContent,
		Message: message,
	})
}
`

// ErrorResponseTemplate - Template for generating ErrorResponse
const ErrorResponseTemplate = `package response

import (
	"fmt" // Added for the InternalError function
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int         ` + "`" + `json:"code"` + "`" + `
	Message string      ` + "`" + `json:"message"` + "`" + `
	Errors  interface{} ` + "`" + `json:"errors,omitempty"` + "`" + `
}

func Error(c *gin.Context, code int, message string, errors interface{}) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
		Errors:  errors,
	})
}

func BadRequest(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusBadRequest, message, errors)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func InternalError(c *gin.Context, err error) {
	Error(c, http.StatusInternalServerError, "Internal Server Error", fmt.Sprintf("error: %v", err))
}

func ValidationError(c *gin.Context, errs map[string]string) {
	Error(c, http.StatusUnprocessableEntity, "Validation failed", errs)
}

func NotFoundResponse(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}
`

// PaginationResponseTemplate - Template for generating pagination response
const PaginationResponseTemplate = `package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page       int         ` + "`" + `json:"page"` + "`" + `
	Limit      int         ` + "`" + `json:"limit"` + "`" + `
	TotalRows  int64       ` + "`" + `json:"total_rows"` + "`" + `
	TotalPages int         ` + "`" + `json:"total_pages"` + "`" + `
	Data       interface{} ` + "`" + `json:"data"` + "`" + `
}

func Paginated(c *gin.Context, data interface{}, page, limit int, totalRows int64) {
	totalPages := int((totalRows + int64(limit) - 1) / int64(limit)) // ceil

	c.JSON(http.StatusOK, Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       data,
	})
}
`
