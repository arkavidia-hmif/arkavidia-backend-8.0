package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	authConfig "arkavidia-backend-8.0/competition/config/authentication"
	"arkavidia-backend-8.0/competition/repository"
)

type AuthRole string

const (
	Admin AuthRole = "Admin"
	Team  AuthRole = "Team"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	ID   uint     `json:"id"`
	Role AuthRole `json:"role"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := authConfig.Config.GetMetadata()
		response := repository.Response[string]{}

		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response.Message = "ERROR: NO TOKEN PROVIDED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
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
			response.Message = "ERROR: TOKEN CANNOT BE PARSED"
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
		if !authToken.Valid {
			response.Message = "ERROR: CLAIMS INVALID"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		c.Set("id", authClaim.ID)
		c.Set("role", authClaim.Role)
		c.Next()
	}
}
