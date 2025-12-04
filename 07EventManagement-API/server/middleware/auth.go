package middleware

import (
	"strings"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var myKey = []byte(config.AppConfig.JWTKEY)

// AuthMiddleware validates Bearer token and sets userId & role in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(400, gin.H{"msg": "No token provided"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(400, gin.H{"msg": "invalid token format"})
			c.Abort()
			return
		}

		myToken := parts[1]

		// Parse with MapClaims so we can inspect fields
		token, err := jwt.ParseWithClaims(myToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			return myKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"msg": "Invalid or Expired Token❌"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(400, gin.H{"msg": "No Data found in token"})
			c.Abort()
			return
		}

		// Try both "id" (preferred) and "userId" (compat)
		var userStrId string
		if v, exists := claims["id"]; exists {
			if s, ok := v.(string); ok {
				userStrId = s
			}
		}
		if userStrId == "" {
			if v, exists := claims["userId"]; exists {
				if s, ok := v.(string); ok {
					userStrId = s
				}
			}
		}

		if userStrId == "" {
			c.JSON(400, gin.H{"msg": "No userId Data found in token"})
			c.Abort()
			return
		}

		userId, err := primitive.ObjectIDFromHex(userStrId)
		if err != nil {
			c.JSON(400, gin.H{"msg": "error converting to user id or expired token"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			c.JSON(400, gin.H{"msg": "No role Data found in token"})
			c.Abort()
			return
		}

		// set the role and userid in context variable
		c.Set("userId", userId)
		c.Set("role", role)

		c.Next()
	}
}

// GenerateJWT creates a token using only email (legacy/compat)
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  "user",
		"exp":   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"iat":   jwt.NewNumericDate(time.Now()),
	})
	return token.SignedString(myKey)
}

// GenerateOAuthJWT creates JWT containing the user's Mongo DB ID (use this for OAuth flows)
func GenerateOAuthJWT(userID, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID, // keep key "id" so middleware picks it up
		"email": email,
		"role":  role,
		"iat":   jwt.NewNumericDate(time.Now()),
		"exp":   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(myKey)
}
