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
	posts, err := p.Document.GetPublicContent(context.Background())
	if err != nil {
		web.AbortHttp(c, err)
		return
	}

	var ids []any
	for _, post := range posts {
		ids = append(ids, post.Id)
	}

	if err := doc.NewLike().CountLike(context.Background(), &posts, ids...); err != nil {
		web.AbortHttp(c, err)
		return
	}

	if err := doc.NewComment().CountComments(context.Background(), &posts, ids...); err != nil {
		web.AbortHttp(c, err)
		return
	}

	for i := 0; i < len(posts); i++ {
		posts[i].Text = h.Decryption(posts[i].Text)
	}

	web.WriteResponse(c, web.WebResponse{
		Code:    200,
		Message: "OK",
		Data:    posts,
	})
}
