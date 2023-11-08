package routes

import (
	c "github.com/forumGamers/post-service-read/controllers"
	md "github.com/forumGamers/post-service-read/middlewares"
	"github.com/gin-gonic/gin"
)

func (r routes) commentRoutes(rg *gin.RouterGroup, cc c.CommentController) {
	uri := rg.Group("/comments")

	// uri.Use(md.Authentication)
	uri.GET("/:postId", cc.FindCommentByPostId)
}
