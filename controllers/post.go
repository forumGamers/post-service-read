package controllers

import (
	"context"

	doc "github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/web"
	"github.com/gin-gonic/gin"
)

type PostController interface {
	FindById(c *gin.Context)
	PublicContent(c *gin.Context)
}

type PostControllerImpl struct {
	Document doc.PostService
}

func NewPostController(db doc.PostService) PostController {
	return &PostControllerImpl{
		Document: db,
	}
}

func (p *PostControllerImpl) FindById(c *gin.Context) {
	get, err := p.Document.FindById(context.Background(), c.Param("postId"))
	if err != nil {
		web.AbortHttp(c, err)
		return
	}

	var post doc.PostDocument
	if err := h.JsonToStruct(get, &post); err != nil {
		web.AbortHttp(c, err)
		return
	}

	post.Text = h.Decryption(post.Text)

	web.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "Success",
		Data:    post,
	})
}

func (p *PostControllerImpl) PublicContent(c *gin.Context) {
	
}