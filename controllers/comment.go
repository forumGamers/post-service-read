package controllers

import (
	"context"

	doc "github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/web"
	"github.com/gin-gonic/gin"
)

type CommentController interface {
	FindCommentByPostId(c *gin.Context)
}

type CommentControllerImpl struct {
	Document doc.CommentService
}

func NewCommentController(db doc.CommentService) CommentController {
	return &CommentControllerImpl{
		Document: db,
	}
}

func (com *CommentControllerImpl) FindCommentByPostId(c *gin.Context) {
	comments, total, err := com.Document.FindCommentByPostId(context.Background(), c.Param("postId"))
	if err != nil {
		web.AbortHttp(c, h.ElasticError(err))
		return
	}

	web.WriteResponseWithMetadata(c, web.WebResponse{
		Code:    200,
		Message: "OK",
		Data:    comments,
	}, web.MetaData{
		Limit:    10,
		Relation: total.Relation,
		Total:    total.Total,
		Page:     comments[len(comments)-1].SearchAfter,
	})

}
