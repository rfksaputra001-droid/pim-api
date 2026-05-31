package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims := &Claims{}
		parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa"})
			return
		}

		c.Set("admin", claims)
		c.Next()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("admin").(*Claims)
		if claims.Role != "SUPER_ADMIN" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak. Hanya Super Admin."})
			return
		}
		c.Next()
	}
}
