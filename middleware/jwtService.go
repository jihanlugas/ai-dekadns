package middleware

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

// JWTService is a contract of what jwtService can do
type JWTService interface {
	ValidateToken(token string) (*jwt.Token, error)
	GetUserByParam(authHeader string, params string) string
}

type jwtService struct {
	secretKey string
	issuer    string
}

// NewJWTService method is creates a new instance of JWTService
func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "secret",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("API_SECRET")
	if secretKey != "" {
		secretKey = "lintasartacl4ud2020"
	}
	return secretKey
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
}

func (j *jwtService) GetUserByParam(authHeader string, params string) string {
	token, _ := j.ValidateToken(authHeader)
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		result := fmt.Sprintf("%s", claims[params])
		return result
	}
	return ""
}
