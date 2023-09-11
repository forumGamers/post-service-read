package routes

import (
	"os"

	db "github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	md "github.com/forumGamers/post-service-read/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type routes struct {
	router *gin.Engine
}

func NewRoutes() {
	h.PanicIfError(godotenv.Load())

	r := routes{router: gin.Default()}

	// groupRoutes := r.router.Group("/api")

	r.router.GET("/ping", func(c *gin.Context) {
		info, code, err := db.Ping()

		c.JSON(200, gin.H{
			"Info":  info,
			"Code":  code,
			"Error": err,
		})
	})

	r.router.Use(md.CheckOrigin)
	r.router.Use(md.Cors())

	port := os.Getenv("PORT")

	if port == "" {
		port = "4301"
	}

	r.router.Run(":" + port)
}
