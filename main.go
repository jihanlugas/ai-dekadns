package main

import (
	"ai-dekadns/database"
	"ai-dekadns/middleware"
	"ai-dekadns/router"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error

	//err = godotenv.Load()
	//if err != nil {
	//	log.Fatalf("Error getting env, %v", err)
	//}

	db := database.GetCorePostsqlConn()
	elastic := database.GetElasticConn()

	dbEvent := os.Getenv("DB_EVENT")
	if dbEvent == "rollback_migrate" || dbEvent == "rollback" {
		database.RunRollback()
	}
	if dbEvent == "migrate_only" {
		database.RunMigration()
	}
	if dbEvent == "migrate" || dbEvent == "rollback_migrate" {
		database.RunMigration()
		database.SeedTypes()
	}

	r := gin.Default()

	r.Use(middleware.LoggerToElastic(elastic))
	r.Use(CORSMiddleware())

	router.Setup(r, db)

	// server
	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Println(err)
	}

}

// CORSMiddleware ..
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		c.Writer.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		c.Writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			//c.Next()
			return
		}
		c.Next()
	}
}
