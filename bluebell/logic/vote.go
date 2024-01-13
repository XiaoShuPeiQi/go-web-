package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"github.com/gin-gonic/gin"
)

func VoteForPost(ctx *gin.Context, p *models.ParamVoteForPost) error {
	useridAny, _ := ctx.Get(JwtUserIDKey)
	userid := useridAny.(int64)
	post := &models.VoteForPost{
		ParamVoteForPost: p,
		UserID:           userid,
	}

	return redis.VoteForPost(post)
}
