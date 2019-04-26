package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if len(tokenString) == 0 {
			context.JSON(http.StatusUnauthorized, "unauthorized")
			context.Abort()
			return
		}

		// get the token from the cookie
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(GetConfig().GetString("jwt_secret")), nil
		})

		// check for valid cookie token
		if err == nil && token.Valid {
			// set the context header
			context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			context.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

			context.Next()
		} else {
			context.JSON(http.StatusUnauthorized, "unauthorized")
			context.Abort()
		}
	}
}
