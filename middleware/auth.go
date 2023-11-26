package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr, err := ctx.Cookie("Authorization")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// decode jwt token
		token, tokenErr := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if tokenErr != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// check if jwt token expired
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				ctx.AbortWithStatus(http.StatusUnauthorized)
			}
			// retrieve and set user in session
			user := authService.GetUserByID(claims["sub"].(uint))
			if user == nil {
				ctx.AbortWithStatus(http.StatusUnauthorized)
			}
			ctx.Set("user", user)
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
