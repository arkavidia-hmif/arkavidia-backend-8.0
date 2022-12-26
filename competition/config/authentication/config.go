package authentication

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthMetadata struct {
	ApplicationName         string
	LoginExpirationDuration time.Duration
	JWTSigningMethod        jwt.SigningMethod
	JWTSignatureKey         []byte
}

type AuthConfig struct {
	metadata AuthMetadata
	once     sync.Once
}

// Private
func (authConfig *AuthConfig) lazyInit() {
	authConfig.once.Do(func() {
		applicationName := os.Getenv("APPLICATION_NAME")
		numberOfSeconds, err := strconv.Atoi(os.Getenv("LOGIN_EXPIRATION_DURATION"))
		if err != nil {
			panic(err)
		}
		loginExpirationDuration := time.Duration(numberOfSeconds) * time.Second
		jwtSigningMethod := jwt.SigningMethodHS256
		jwtSignatureKey := []byte(os.Getenv("JWT_SIGNATURE_KEY"))

		authConfig.metadata.ApplicationName = applicationName
		authConfig.metadata.LoginExpirationDuration = loginExpirationDuration
		authConfig.metadata.JWTSigningMethod = jwtSigningMethod
		authConfig.metadata.JWTSignatureKey = jwtSignatureKey
	})
}

// Public
func (authConfig *AuthConfig) GetMetadata() AuthMetadata {
	authConfig.lazyInit()
	return authConfig.metadata
}

var Config = &AuthConfig{}
