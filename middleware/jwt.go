package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const NoHeaderProvidedErr = "No Authorization header provided"

func GetTokenRemainingValidity(timestamp interface{}) bool {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := time.Until(tm)
		if remainer > 0 {
			return true
		}
	}
	return false
}

// AuthorizeBasic validates the token user given, return 401 if not valid
func AuthorizeBasic() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(403, gin.H{"error": NoHeaderProvidedErr})
			c.Abort()
			return
		}

		apiKey := os.Getenv("API_SECRET")
		if authHeader != apiKey {
			c.JSON(403, gin.H{"error": "Token is not valid", "des": "Not valid token"})
			c.Abort()
			return
		}
	}
}

// AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(jwtService JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(403, gin.H{"error": NoHeaderProvidedErr})
			c.Abort()
			return
		}
		extractedToken := strings.Split(authHeader, "Bearer ")
		authHeader = strings.TrimSpace(extractedToken[1])
		fmt.Println("Authorization", authHeader)

		token, err := jwtService.ValidateToken(authHeader)
		if !token.Valid {
			c.JSON(403, gin.H{"error": "Token is not valid", "des": err.Error()})
			c.Abort()
			return
		}
	}
}

func HasJWT(jwtService JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(403, gin.H{"error": NoHeaderProvidedErr})
			c.Abort()
			return
		}
		extractedToken := strings.Split(authHeader, "Bearer ")
		authHeader = strings.TrimSpace(extractedToken[1])
		fmt.Println("Authorization", authHeader)

		_, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			c.JSON(403, gin.H{"error": "Token error", "des": err.Error()})
			c.Abort()
			return
		}
	}
}

func AuthorizeSuperAdminMiddleware(jwtService JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractTokenFromHeader(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not Valid Token or Token Expired"})
			return
		}

		claims, err := getClaimsFromToken(token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err := validateTokenClaims(ctx, claims); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("userID", claims["user_id"])
		ctx.Set("role", claims["role"])
		ctx.Set("roleID", claims["roleID"])
	}
}

func extractTokenFromHeader(ctx *gin.Context) (string, error) {
	const BearerSchema = "bearer "
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New(NoHeaderProvidedErr)
	}
	return authHeader[len(BearerSchema):], nil
}

func getClaimsFromToken(token *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func validateTokenClaims(ctx *gin.Context, claims jwt.MapClaims) error {
	if !GetTokenRemainingValidity(claims["exp"]) {
		return errors.New("Token Expired")
	}

	if claims["nv"] == true {
		return errors.New("Validate Device")
	}

	if claims["role"] != "Superadmin" {
		return errors.New("Permission denied")
	}

	return nil
}

func AuthorizeMiddlewareWithAddContext(jwtService JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractTokenFromHeader(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not Valid Token or Token Expired"})
			return
		}

		claims, err := getClaimsFromToken(token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err := validateTokenClaimsAddCtx(ctx, claims); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set user data to context
		setContextWithClaims(ctx, claims)
	}
}

func validateTokenClaimsAddCtx(ctx *gin.Context, claims jwt.MapClaims) error {
	if !GetTokenRemainingValidity(claims["exp"]) {
		return errors.New("Token Expired")
	}

	if claims["nv"] == true {
		return errors.New("Validate Device")
	}

	return nil
}

func setContextWithClaims(ctx *gin.Context, claims jwt.MapClaims) {
	ctx.Set("userID", claims["user_id"])
	ctx.Set("role", claims["role"])
	ctx.Set("roleID", claims["roleID"])
	ctx.Set("organizationID", claims["organization_id"])
}

// api basic auth
func AuthorizeCronMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("x-auth-cron")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": NoHeaderProvidedErr})
			return

		}

		cronSecretKey := os.Getenv("API_SECRET")
		if authHeader != cronSecretKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Not Valid Token"})
			return

		}
	}
}
