package middlewares

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/koki/randommatch/utils/token"
)

func extract(r *http.Request) string {
	authValues := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authValues) == 2 && authValues[0] == "Bearer" {
		return authValues[1]
	}
	return ""
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := extract(c.Request)
		claims, err := token.Validate(bearerToken)
		if err != nil {
			log.Println("bearer token error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Operation denied"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("tokenClaims", claims)

		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Location"},
		MaxAge:        12 * time.Hour,
	})
}
