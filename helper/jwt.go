package helper

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type jwtCustomClaim struct {
	UserID         string `json:"user_id"`
	ProjectID      string `json:"project_id"`
	OrganizationID string `json:"organization_id"`

	jwt.StandardClaims
}

type refCustomClaim struct {
	UserID         string `json:"user_id"`
	ProjectID      string `json:"project_id"`
	OrganizationID string `json:"organization_id"`
	jwt.StandardClaims
}

type JwtService interface {
	GenerateToken(UserID string, ProjectID string, OrganizationID string) string
	RefreshToken(UserID string, ProjectID string, OrganizationID string) string
	ValidateToken(signedToken string) (claims *jwtCustomClaim, err error)
}

type jwtService struct {
	secretKey string
	issure    string
}

func NewJwtService() JwtService {
	return &jwtService{secretKey: os.Getenv("JWT_SECRET"), issure: "cloudsecret"}
}

func (s *jwtService) GenerateToken(UserID string, ProjectID string, OrganizationID string) string {
	claims := &jwtCustomClaim{
		UserID, ProjectID, OrganizationID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 180).Unix(),
			Issuer:    s.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (s *jwtService) RefreshToken(UserID string, ProjectID string, OrganizationID string) string {
	claims := &refCustomClaim{
		UserID, ProjectID, OrganizationID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 200).Unix(),
			Issuer:    s.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (s *jwtService) ValidateToken(signedToken string) (claims *jwtCustomClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtCustomClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secretKey), nil
		},
	)

	if err != nil {
		return
	}
	claims, ok := token.Claims.(*jwtCustomClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ParseToken(tokenString string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	} else {
		return nil
	}
}

func GetUserID(c *gin.Context) string {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string)
	} else {
		return ""
	}
}

func GetRole(c *gin.Context) string {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["role"].(string)
	} else {
		return ""
	}
}

func GetOrganizationID(c *gin.Context) string {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["organization_id"].(string)
	} else {
		return ""
	}
}

func GetIP(c *gin.Context) string {
	//Get IP from the X-REAL-IP header
	ip := c.Request.Header.Get("X-Real-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	//Get IP from X-FORWARDED-FOR header
	ips := c.Request.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return ""
}
