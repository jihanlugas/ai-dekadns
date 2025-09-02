package router

import (
	"ai-dekadns/app/organization"
	"ai-dekadns/app/project"
	"ai-dekadns/app/record"
	"ai-dekadns/app/superadminrole"
	"ai-dekadns/app/types"
	"ai-dekadns/app/user"
	"ai-dekadns/app/zone"
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
	organizationRepository := organization.NewRepository(db)
	projectRepository := project.NewRepository(db)
	userRepository := user.NewRepository(db)
	superadminroleRepository := superadminrole.NewRepository(db)
	zoneRepository := zone.NewRepository(db)
	typesRepository := types.NewRepository(db)

	// Usecase
	zoneUsecase := zone.NewUsecase(organizationRepository, projectRepository, userRepository, superadminroleRepository, zoneRepository)
	recordUsecase := record.NewUsecase(organizationRepository, projectRepository, userRepository, superadminroleRepository, zoneRepository)
	typesUsecase := types.NewUsecase(typesRepository)

	// Handler
	zoneHandler := zone.NewHandler(zoneUsecase)
	recordHandler := record.NewHandler(recordUsecase)
	typesHandler := types.NewHandler(typesUsecase)

	router := c

	router.GET("/check", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	})

	zoneRouter := router.Group("/zone")
	zoneRouter.GET("", middleware.AuthorizeJWT(jwtService), zoneHandler.Page)
	zoneRouter.POST("", middleware.AuthorizeJWT(jwtService), zoneHandler.Create)
	zoneRouter.GET("/:id", middleware.AuthorizeJWT(jwtService), zoneHandler.GetById)
	zoneRouter.DELETE("/:id", middleware.AuthorizeJWT(jwtService), zoneHandler.Delete)

	recordRouter := router.Group("/record")
	recordRouter.POST("", middleware.AuthorizeJWT(jwtService), recordHandler.Create)
	recordRouter.PUT("", middleware.AuthorizeJWT(jwtService), recordHandler.Update)
	recordRouter.DELETE("", middleware.AuthorizeJWT(jwtService), recordHandler.Delete)

	typesRouter := router.Group("/types")
	typesRouter.GET("", middleware.AuthorizeJWT(jwtService), typesHandler.Page)

}
