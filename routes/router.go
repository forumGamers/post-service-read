package routes

import (
	"os"

	"github.com/forumGamers/post-service-read/controllers"
	db "github.com/forumGamers/post-service-read/database"
	md "github.com/forumGamers/post-service-read/middlewares"
	"github.com/gin-gonic/gin"
)

type routes struct {
	router *gin.Engine
}

func NewRoutes(post controllers.PostController, comment controllers.CommentController) {
	r := routes{router: gin.Default()}

	groupRoutes := r.router.Group("/api/v1")

	r.router.Use(md.CheckOrigin)
	r.router.Use(md.Cors())
	r.router.GET("/ping", func(c *gin.Context) {
		info, code, err := db.Ping()

		c.JSON(200, gin.H{
			"Info":  info,
			"Code":  code,
			"Error": err,
		})
	})
	r.postRoute(groupRoutes, post)
	r.commentRoutes(groupRoutes, comment)

	port := os.Getenv("PORT")

	if port == "" {
		port = "4301"
	}

	r.router.Run(":" + port)
}
