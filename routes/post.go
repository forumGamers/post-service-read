package routes

import (
	"github.com/forumGamers/post-service-read/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) postRoute(rg *gin.RouterGroup, pc controllers.PostController) {
	uri := rg.Group("/post")

	uri.GET("/:postId", pc.FindById)
}
