package routes

import (
	"github.com/forumGamers/post-service-read/controllers"
	md "github.com/forumGamers/post-service-read/middlewares"
	"github.com/gin-gonic/gin"
)

func (r routes) postRoute(rg *gin.RouterGroup, pc controllers.PostController) {
	uri := rg.Group("/post")

	uri.Use(md.Authentication)
	uri.GET("/public", pc.PublicContent)
	uri.GET("/:postId", pc.FindById)
}
