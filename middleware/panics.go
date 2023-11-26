package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var errorMsg string
				var statusCode int
				switch v := r.(type) {
				case error:
					errorMsg = v.Error()
					statusCode = http.StatusInternalServerError
				case string:
					errorMsg = v
					statusCode = http.StatusBadRequest
					return
				default:
					errorMsg = "Unknown error occurred"
					statusCode = http.StatusInternalServerError
				}
				ctx.JSON(statusCode, gin.H{
					"error": errorMsg,
				})
			}
		}()
		ctx.Next()
	}
}
