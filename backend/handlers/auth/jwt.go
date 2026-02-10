package auth

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (any, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			log.Printf("error while creating token: %v\n", err)
			return
		}

		userID, _ := strconv.Atoi(claims.Subject)
		c.Set("userId", userID)
		c.Next()
	}
}

func GenerateToken(userId int) (string, error) {
	claims := jwt.RegisteredClaims {
		ExpiresAt:	jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		IssuedAt:	jwt.NewNumericDate(time.Now()),
		Subject:	strconv.Itoa(userId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}