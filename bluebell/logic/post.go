package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const JwtUserIDKey = "userid"

/*
CreatePost
创建一个帖子
*/
func CreatePost(ctx *gin.Context, p *models.Post) (err error) {
	//生成post_id
	p.PostID = snowflake.GetID()
	//设置用户id
	value, exists := ctx.Get(JwtUserIDKey)
	if exists {
		p.AuthorID = value.(int64)
	}
	// 调用dao层完成插入工作
	if err = mysql.InsertPost(p); err != nil {
		return
	}
	//写入redis的发帖时间zset中
	return redis.CreatePost(p.PostID)
}

/*
GetPost
根据帖子id查询帖子信息
*/
func GetPost(id int64) (post *models.PostDetail, err error) {
	post = new(models.PostDetail)

	//	1.根据postid查内容
	post, err = mysql.GetPostByID(id)
	if err != nil {
		return
	}
	//	2.根据用户id查姓名
	var username string
	username, err = mysql.GetUserNameByID(post.AuthorID)
	if err != nil {
		return
	}
	post.AuthorName = username
	//	3.根据社区id查姓名，详情
	var name string
	var intro string
	name, intro, err = mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		return
	}
	post.CommuName = name
	post.Introduction = intro
	//	4.整合数据
	return
}

/*
GetPostList
返回帖子列表
*/
func GetPostList(size int64, page int64) (data []*models.PostDetail, err error) {
	//	调用dao层得到帖子post列表
	pList, err := mysql.GetPostList(size, page)
	if err != nil {
		zap.L().Debug("GetPostList调试信息", zap.Int("plist:", len(pList)))
		return nil, err
	}

	//	循环添加详细信息
	for _, v := range pList {
		name, err := mysql.GetUserNameByID(v.AuthorID)
		if err != nil {
			return nil, err
		}
		cName, intro, err := mysql.GetCommunityByID(v.CommunityID)
		if err != nil {
			return nil, err
		}
		OneData := &models.PostDetail{
			Post:         v,
			AuthorName:   name,
			CommuName:    cName,
			Introduction: intro,
		}
		data = append(data, OneData)
	}
	return data, nil
}
func GetPostList2(p *models.ParamPostList) ([]*models.PostDetail, error) {
	//1.去redis层查有序的id列表
	ids, err := redis.GetPostIDsByOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("ids 为空！")
		return nil, err
	}
	//查投赞成票信息
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	//2.用id去dao层查帖子信息
	pList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error(" mysql.GetPostListByIDs wrong", zap.Error(err))
		return nil, err
	}
	data := make([]*models.PostDetail, 0)
	//3.每个帖子完善作者信息和社区信息
	for i, v := range pList {
		name, err := mysql.GetUserNameByID(v.AuthorID)
		if err != nil {
			return nil, err
		}
		cName, intro, err := mysql.GetCommunityByID(v.CommunityID)
		if err != nil {
			return nil, err
		}
		OneData := &models.PostDetail{
			Post:         v,
			VoteScore:    voteData[i],
			AuthorName:   name,
			CommuName:    cName,
			Introduction: intro,
		}
		data = append(data, OneData)
	}
	return data, nil
}
