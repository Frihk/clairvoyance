package utils

import "github.com/gin-gonic/gin"

// SuccessResponse writes a standard JSON success envelope with the provided
// status and data payload.
func SuccessResponse(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

// ErrorResponse writes a standard JSON error envelope with the provided status
// and error message.
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}

// PaginatedResponse writes a paginated JSON response including meta fields.
func PaginatedResponse(c *gin.Context, status int, data any, total int64, page, limit int) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
