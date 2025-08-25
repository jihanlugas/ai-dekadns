package router

import (
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
	//sslRepository := ssl.NewRepository(coreDb, sslDb)
	//organizationRepository := organization.NewRepository(coreDb, sslDb)
	//projectRepository := project.NewRepository(coreDb, sslDb)
	//userRepository := user.NewRepository(coreDb, sslDb)
	//superadminroleRepository := superadminrole.NewRepository(coreDb, sslDb)

	//// Usecase
	//sslUsecase := ssl.NewUsecase(sslRepository, organizationRepository, projectRepository, userRepository, superadminroleRepository)
	//
	//// Handler
	//sslHandler := ssl.NewHandler(sslUsecase)

	router := c

	router.GET("/check", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})

	//router.GET("/", middleware.AuthorizeJWT(jwtService), sslHandler.Page)
	//router.POST("/import", middleware.AuthorizeJWT(jwtService), sslHandler.ImportSSL)

}
