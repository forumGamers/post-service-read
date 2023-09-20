package controllers

import (
	"context"
	"sync"

	doc "github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	i "github.com/forumGamers/post-service-read/interfaces"
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

	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(3)

	go func(posts []i.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- doc.NewLike().CountLike(context.Background(), &posts, "", ids...)
	}(posts, ids...)

	go func(posts []i.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- doc.NewComment().CountComments(context.Background(), &posts, ids...)
	}(posts, ids...)

	go func(posts []i.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- doc.NewShare().CountShares(context.Background(), &posts, "", ids...)
	}(posts, ids...)

	var errors error
	flag := false
	for i := 0; i < 3; i++ {
		select {
		case err := <-errCh:
			{
				if err != nil && !flag {
					flag = true
					errors = err
				}
			}
		}
	}

	wg.Wait()
	if flag {
		web.AbortHttp(c, errors)
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
