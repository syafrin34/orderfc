package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing required auth",
			})
			c.Abort()
			return
		}
		// format token
		// // Authorization : Bearer xx
		// Bearer [xxx]
		tokenString := strings.Split(authHeader, " ")
		fmt.Println("token ", tokenString)
		if len(tokenString) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token 1",
			})
			c.Abort()
			return
		}

		// token, err := jwt.Parse(tokenString[1], func(t *jwt.Token) (interface{}, error) {
		// 	return []byte(secret), nil
		// })

		token, err := jwt.Parse(tokenString[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		fmt.Println("token valid ", token.Valid)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token 2",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		fmt.Println("token ke 3 ", claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token 3",
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims["user_id"].(float64))
		c.Next()

	}

}
