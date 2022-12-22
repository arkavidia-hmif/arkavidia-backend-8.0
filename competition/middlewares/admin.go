package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	authConfig "arkavidia-backend-8.0/competition/config/authentication"
)

type AdminClaims struct {
	jwt.RegisteredClaims
	AdminID uint `json:"admin_id"`
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := authConfig.GetAuthConfig()

		adminHeader := c.GetHeader("Authorization")
		if !strings.Contains(adminHeader, "Bearer") {
			response := gin.H{"Message": "ERROR: NO TOKEN PROVIDED"}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		adminString := strings.Replace(adminHeader, "Bearer ", "", -1)
		adminClaim := AdminClaims{}
		adminToken, err := jwt.ParseWithClaims(adminString, &adminClaim, func(adminToken *jwt.Token) (interface{}, error) {
			if method, ok := adminToken.Method.(*jwt.SigningMethodHMAC); !ok || method != config.JWTSigningMethod {
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
		if !adminToken.Valid {
			response := gin.H{"Message": "CLAIMS INVALID"}
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		c.Set("admin_id", adminClaim.AdminID)
		c.Next()
	}
}
