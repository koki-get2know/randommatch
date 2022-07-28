package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/koki/randommatch/utils/token"
)

func extract(c *gin.Context) string {
	authValues := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(authValues) == 2 && authValues[0] == "Bearer" {
		return authValues[1]
	}
	return ""
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := extract(c)

		if _, err := token.Validate(bearerToken); err != nil {
			fmt.Println("bearer token error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Operation denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		MaxAge:       12 * time.Hour,
	})
}
