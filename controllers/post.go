package controllers

import (
	"context"
	"sync"
	"time"

	doc "github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/comment"
	"github.com/forumGamers/post-service-read/pkg/like"
	"github.com/forumGamers/post-service-read/pkg/post"
	v "github.com/forumGamers/post-service-read/validator"
	"github.com/forumGamers/post-service-read/web"
	"github.com/gin-gonic/gin"
)

type PostController interface {
	FindById(c *gin.Context)
	PublicContent(c *gin.Context)
	FindMyPost(c *gin.Context)
}

type PostControllerImpl struct {
	Document post.PostService
}

func NewPostController(db post.PostService) PostController {
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

	var post post.PostDocument
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
	var query web.PostParams
	c.ShouldBindQuery(&query)
	v.ValidatePostQuery(&query)

	posts, total, err := p.Document.GetPublicContent(context.Background(), query)
	if err != nil {
		web.AbortHttp(c, h.ElasticError(err))
		return
	}

	var ids []any
	for _, post := range posts {
		ids = append(ids, post.Id)
	}

	errCh := make(chan error)
	uuid := doc.GetUser(c).UUID
	var wg sync.WaitGroup
	wg.Add(3)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- like.NewLike().CountLike(context.Background(), &posts, uuid, ids...)
	}(posts, ids...)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- comment.NewComment().CountComments(context.Background(), &posts, ids...)
	}(posts, ids...)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond) //agar ga error 429
		errCh <- doc.NewShare().CountShares(context.Background(), &posts, uuid, ids...)
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
		web.AbortHttp(c, h.ElasticError(errors))
		return
	}

	web.WriteResponseWithMetadata(c, web.WebResponse{
		Code:    200,
		Message: "OK",
		Data:    posts,
	}, web.MetaData{
		Limit:    query.Limit,
		Relation: total.Relation,
		Total:    total.Total,
		Page:     posts[len(posts)-1].SearchAfter,
	})
}

func (p *PostControllerImpl) FindMyPost(c *gin.Context) {
	// posts, total, err := p.Document.FindByUserId(context.Background(), doc.GetUser(c).UUID)
	// if err != nil {
	// 	web.AbortHttp(c, h.ElasticError(err))
	// 	return
	// }

}
