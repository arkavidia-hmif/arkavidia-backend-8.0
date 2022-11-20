package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	authConfig "arkavidia-backend-8.0/competition/config/authentication"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	TeamID uint `json:"team_id"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := authConfig.GetAuthConfig()

		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := gin.H{"Message": "ERROR: NO TOKEN PROVIDED"}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		authString := strings.Replace(authHeader, "Bearer ", "", -1)
		authClaim := AuthClaims{}
		authToken, err := jwt.ParseWithClaims(authString, &authClaim, func(authToken *jwt.Token) (interface{}, error) {
			if method, ok := authToken.Method.(*jwt.SigningMethodHMAC); !ok || method != config.JWTSigningMethod {
				return nil, fmt.Errorf("ERROR: SIGNING METHOD INVALID")
			}
			return config.JWTSignatureKey, nil
		})
		if err != nil {
			response := gin.H{"Message": err.Error()}
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}
		if !authToken.Valid {
			response := gin.H{"Message": "CLAIMS INVALID"}
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		c.Set("team_id", authClaim.TeamID)
		c.Next()
	}
}
