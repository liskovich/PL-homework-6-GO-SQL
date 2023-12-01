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
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
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
			userIDflt := claims["sub"]
			floatValue, ok := userIDflt.(float64)
			var userID uint
			if !ok {
				ctx.AbortWithStatus(http.StatusUnauthorized)
			} else {
				userID = uint(floatValue)
			}

			user := authService.GetUserByID(userID)
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

func UserDetailMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr, err := ctx.Cookie("Authorization")
		if err != nil {
			ctx.Next()
			return
		}
		// otherwise, try to get logged in user details
		// if something below fails, assume that the user is not logged in and move on

		// decode jwt token
		token, tokenErr := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if tokenErr != nil {
			ctx.Next()
			return
		}

		// TODO: fix error - runtime error: invalid memory address or nil pointer dereference
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// check if jwt token expired
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				ctx.Next()
			}

			// retrieve and set user in session
			userIDflt := claims["sub"]
			floatValue, ok := userIDflt.(float64)
			var userID uint
			if !ok {
				ctx.Next()
				return
			} else {
				userID = uint(floatValue)
			}

			user := authService.GetUserByID(userID)
			if user == nil {
				ctx.Next()
				return
			}
			ctx.Set("user", user)
			ctx.Next()
		} else {
			ctx.Next()
			return
		}
	}
}
