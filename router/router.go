package router

import (
	"ai-dekadns/app/dns"
	"ai-dekadns/app/organization"
	"ai-dekadns/app/project"
	"ai-dekadns/app/superadminrole"
	"ai-dekadns/app/user"
	"ai-dekadns/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var jwtService middleware.JWTService

func init() {
	jwtService = middleware.NewJWTService()
}

func Setup(c *gin.Engine, db *gorm.DB) {
	// Repository
	dnsRepository := dns.NewRepository(db)
	organizationRepository := organization.NewRepository(db)
	projectRepository := project.NewRepository(db)
	userRepository := user.NewRepository(db)
	superadminroleRepository := superadminrole.NewRepository(db)

	// Usecase
	dnsUsecase := dns.NewUsecase(dnsRepository, organizationRepository, projectRepository, userRepository, superadminroleRepository)

	// Handler
	dnsHandler := dns.NewHandler(dnsUsecase)

	router := c

	router.GET("/check", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})

	router.POST("", middleware.AuthorizeJWT(jwtService), dnsHandler.Create)

}
