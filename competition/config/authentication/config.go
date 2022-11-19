package authentication

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthConfig struct {
	ApplicationName         string
	LoginExpirationDuration time.Duration
	JWTSigningMethod        jwt.SigningMethod
	JWTSignatureKey         []byte
}

var currentAuthConfig *AuthConfig = nil

func Init() *AuthConfig {
	applicationName := os.Getenv("APPLICATION_NAME")
	numberOfSeconds, err := strconv.Atoi(os.Getenv("LOGIN_EXPIRATION_DURATION"))
	if err != nil {
		panic(err)
	}
	loginExpirationDuration := time.Duration(numberOfSeconds) * time.Second
	jwtSigningMethod := jwt.SigningMethodHS256
	jwtSignatureKey := []byte(os.Getenv("JWT_SIGNATURE_KEY"))
	return &AuthConfig{
		ApplicationName:         applicationName,
		LoginExpirationDuration: loginExpirationDuration,
		JWTSigningMethod:        jwtSigningMethod,
		JWTSignatureKey:         jwtSignatureKey,
	}
}

func GetAuthConfig() *AuthConfig {
	if currentAuthConfig == nil {
		currentAuthConfig = Init()
	}
	return currentAuthConfig
}
