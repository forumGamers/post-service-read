package controllers

import (
	"context"
	"sync"
	"time"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/comment"
	"github.com/forumGamers/post-service-read/pkg/like"
	"github.com/forumGamers/post-service-read/pkg/post"
	"github.com/forumGamers/post-service-read/pkg/reply"
	"github.com/forumGamers/post-service-read/pkg/share"
	"github.com/forumGamers/post-service-read/pkg/user"
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
	Document   post.PostService
	LikeDoc    like.LikeService
	CommentDoc comment.CommentService
	ReplyDoc   reply.ReplyService
	ShareDoc   share.ShareService
}

func NewPostController(
	db post.PostService,
	like like.LikeService,
	comment comment.CommentService,
	reply reply.ReplyService,
	share share.ShareService,
) PostController {
	return &PostControllerImpl{
		Document:   db,
		LikeDoc:    like,
		CommentDoc: comment,
		ReplyDoc:   reply,
		ShareDoc:   share,
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
	uuid := user.GetUser(c).UUID
	var wg sync.WaitGroup
	wg.Add(3)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- p.LikeDoc.CountLike(context.Background(), &posts, uuid, ids...)
	}(posts, ids...)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- p.CommentDoc.CountComments(context.Background(), &posts, ids...)
	}(posts, ids...)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond) //agar ga error 429
		errCh <- p.ShareDoc.CountShares(context.Background(), &posts, uuid, ids...)
	}(posts, ids...)

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			web.AbortHttp(c, h.ElasticError(err))
			return
		}
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
	var query web.PostParams
	c.ShouldBindQuery(&query)
	v.ValidatePostQuery(&query)
	uuid := user.GetUser(c).UUID

	posts, total, err := p.Document.FindByUserId(context.Background(), uuid, query)
	if err != nil {
		web.AbortHttp(c, h.ElasticError(err))
		return
	}

	var ids []any
	for _, post := range posts {
		ids = append(ids, post.Id)
	}

	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(3)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- p.LikeDoc.CountLike(context.Background(), &posts, uuid, ids...)
	}(posts, ids...)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		errCh <- p.CommentDoc.CountComments(context.Background(), &posts, ids...)
	}(posts, ids...)

	go func(posts []post.PostResponse, ids ...any) {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond) //agar ga error 429
		errCh <- p.ShareDoc.CountShares(context.Background(), &posts, uuid, ids...)
	}(posts, ids...)

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			web.AbortHttp(c, h.ElasticError(err))
			return
		}
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
