package redis

import (
	"bluebell/models"
	"errors"
	"github.com/go-redis/redis/v8"
	"math"
	"strconv"
	"time"
)

const (
	oneWeek        = float64(time.Hour * 24 * 7)
	oneScore int64 = 432
)

var (
	ErrorOverTime = errors.New("超过可投票时间")
)

// CreatePost 在创建帖子时生成redis分数表
func CreatePost(id int64) error {
	//开启事务，以下两个绑定的
	pipeline := rdb.TxPipeline()
	//插入redis的Zset中
	//	时间表
	pipeline.ZAdd(ctx, KeyTimeZset, &redis.Z{
		//time.Now().Unix()时间戳，自1970年以来的秒数
		Score:  float64(time.Now().Unix()),
		Member: id,
	})
	//	分数表(以时间为起始值)
	pipeline.ZAdd(ctx, KeyScoreZset, &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: id,
	})
	_, err := pipeline.Exec(ctx)
	return err
}

// VoteForPost 投票功能业务逻辑
func VoteForPost(p *models.VoteForPost) error {
	userid := strconv.FormatInt(p.UserID, 10)
	postid := strconv.FormatInt(p.PostID, 10)
	//	1.查询限制条件，只有在1星期里发布的帖子才有资格投票
	postTime, err := rdb.ZScore(ctx, KeyTimeZset, postid).Result()
	if err != nil {
		return err
	}
	if float64(time.Now().Unix())-postTime > oneWeek {
		return ErrorOverTime
	}
	//	2.查询用户先前为该帖子投的什么票 ，并插入分数
	//Zset表：key为帖子id，记录了用户id和投什么票
	oDre := rdb.ZScore(ctx, KeyPostZsetP+postid, userid).Val()

	var i int64 = 1
	if oDre > float64(p.Direction) {
		i = -1
	}
	pipeline := rdb.TxPipeline() //开始插入数据了，两步为同步操作
	abs := math.Abs(float64(p.Direction) - oDre)
	pipeline.ZIncrBy(ctx, KeyScoreZset, float64(i*int64(abs)*oneScore), postid)
	//	3.插入这次的投票纪录
	pipeline.ZAdd(ctx, KeyPostZsetP+postid, &redis.Z{
		Score:  float64(p.Direction),
		Member: userid,
	})
	_, err = pipeline.Exec(ctx)
	return err
}
