package controllers

import (
	"context"
	"sort"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/comment"
	"github.com/forumGamers/post-service-read/pkg/reply"
	v "github.com/forumGamers/post-service-read/validator"
	"github.com/forumGamers/post-service-read/web"
	"github.com/gin-gonic/gin"
)

type CommentController interface {
	FindCommentByPostId(c *gin.Context)
}

type CommentControllerImpl struct {
	Document comment.CommentService
}

func NewCommentController(db comment.CommentService) CommentController {
	return &CommentControllerImpl{
		Document: db,
	}
}

func (com *CommentControllerImpl) FindCommentByPostId(c *gin.Context) {
	var query web.Params
	c.ShouldBind(&query)
	v.ValidateCommentParams(&query)

	comments, total, err := com.Document.FindCommentByPostId(context.Background(), c.Param("postId"), query)
	if err != nil {
		web.AbortHttp(c, h.ElasticError(err))
		return
	}

	var ids []any
	for _, val := range comments {
		ids = append(ids, val.Id)
	}

	replies, err := reply.NewReply().FindCommentsReply(context.Background(), ids)
	if err != nil {
		web.AbortHttp(c, h.ElasticError(err))
		return
	}

	for _, comment := range comments {
		commentReply := make([]reply.Reply, 0)
		for _, reply := range replies {
			if comment.Id == reply.CommentId {
				commentReply = append(commentReply, reply)
			}
		}
		if len(commentReply) > 0 {
			sort.Slice(commentReply, func(i, j int) bool {
				return commentReply[i].CreatedAt.After(commentReply[j].CreatedAt)
			})
		}
		comment.Reply = commentReply
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
